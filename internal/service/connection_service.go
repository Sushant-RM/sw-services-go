package service

import (
	"fmt"
	"log"
	"time"

	"github.com/BrajK111/sw-services/internal/domain/dto"
	"github.com/BrajK111/sw-services/internal/integration"
	"github.com/BrajK111/sw-services/internal/repository/postgres"
	"github.com/BrajK111/sw-services/internal/transport/kafka/producer"
	"github.com/google/uuid"
)

type SewerageService interface {
	CreateConnection(
		reqInfo dto.RequestInfo,
		connection dto.SewerageConnection,
	) (dto.SewerageConnection, error)

	SearchConnection(
		reqInfo dto.RequestInfo,
		tenantID string,
		applicationNumber string,
	) ([]dto.SewerageConnection, error)

	UpdateConnection(
		reqInfo dto.RequestInfo,
		connection dto.SewerageConnection,
	) (dto.SewerageConnection, error)
}

type sewerageService struct {
	repository     postgres.SewerageRepository
	producer       *producer.Producer
	idGenClient    *integration.IDGenClient
	workflowClient *integration.WorkflowClient
}

func NewSewerageService(
	repository postgres.SewerageRepository,
	producer *producer.Producer,
	idGen *integration.IDGenClient,
	workflow *integration.WorkflowClient,
) SewerageService {
	return &sewerageService{
		repository:     repository,
		producer:       producer,
		idGenClient:    idGen,
		workflowClient: workflow,
	}
}

func (s *sewerageService) CreateConnection(
	reqInfo dto.RequestInfo,
	connection dto.SewerageConnection,
) (dto.SewerageConnection, error) {

	id := uuid.New().String()

	// Generate Application Number via IDGen with robust fallback
	var appNo string
	var err error
	if s.idGenClient != nil {
		appNo, err = s.idGenClient.GenerateApplicationNo(reqInfo, connection.TenantID)
	}
	if err != nil || appNo == "" {
		log.Printf("Warning: IDGen unavailable, using fallback UUID for applicationNo: %v", err)
		appNo = fmt.Sprintf("SW-APP-%s", uuid.New().String()[:8])
	}

	connection.ID = id
	connection.ApplicationNo = appNo
	connection.ApplicationNumber = appNo

	// Enrich audit details
	currentTime := time.Now().UnixMilli()
	username := "system"
	if reqInfo.UserInfo != nil {
		username = reqInfo.UserInfo.UserName
	}
	connection.AuditDetails = &dto.AuditDetails{
		CreatedBy:        username,
		LastModifiedBy:   username,
		CreatedTime:      currentTime,
		LastModifiedTime: currentTime,
	}

	if connection.ApplicationStatus == "" {
		connection.ApplicationStatus = "INITIATED"
	}
	connection.Status = connection.ApplicationStatus

	// Save to DB
	err = s.repository.CreateConnection(
		connection,
		appNo,
		id,
	)
	if err != nil {
		return connection, err
	}

	// Trigger workflow process instance if present
	if connection.ProcessInstance.Action == "" {
		connection.ProcessInstance = dto.ProcessInstance{
			Action:          "INITIATE",
			BusinessService: "SW",
			ModuleName:      "SW",
		}
	}

	if s.workflowClient != nil {
		pi, err := s.workflowClient.CallWorkflow(reqInfo, connection)
		if err != nil {
			log.Printf("Warning: Workflow transition failed: %v, using local state machine.", err)
			nextStatus := getLocalNextState("INITIATED", connection.ProcessInstance.Action)
			connection.ApplicationStatus = nextStatus
			connection.Status = nextStatus
		} else if pi != nil {
			// update application status if workflow changed it
			if pi.State != nil && pi.State.ApplicationStatus != "" {
				connection.ApplicationStatus = pi.State.ApplicationStatus
				connection.Status = pi.State.ApplicationStatus
			} else {
				connection.ApplicationStatus = pi.Action
				connection.Status = pi.Action
			}
		}
	} else {
		nextStatus := getLocalNextState("INITIATED", connection.ProcessInstance.Action)
		connection.ApplicationStatus = nextStatus
		connection.Status = nextStatus
	}

	// Publish to Kafka topic save-sw-connection
	kafkaReq := dto.SewerageConnectionRequest{
		RequestInfo:        reqInfo,
		SewerageConnection: connection,
		IsCreateCall:       true,
	}
	_ = s.producer.PublishConnectionCreated(kafkaReq, appNo)

	return connection, nil
}

func (s *sewerageService) SearchConnection(
	reqInfo dto.RequestInfo,
	tenantID string,
	applicationNumber string,
) ([]dto.SewerageConnection, error) {

	connections, err := s.repository.SearchConnection(
		tenantID,
		applicationNumber,
	)
	if err != nil {
		return nil, err
	}

	return connections, nil
}

func (s *sewerageService) UpdateConnection(
	reqInfo dto.RequestInfo,
	connection dto.SewerageConnection,
) (dto.SewerageConnection, error) {

	appNo := connection.ApplicationNo
	if appNo == "" {
		appNo = connection.ApplicationNumber
	}

	// Load existing connection from the database to avoid overwriting unchanged fields with defaults/empty values
	existingConns, err := s.repository.SearchConnection(connection.TenantID, appNo)
	if err != nil {
		return connection, err
	}
	if len(existingConns) == 0 {
		return connection, fmt.Errorf("sewerage connection not found for applicationNo: %s", appNo)
	}

	existing := existingConns[0]

	// Merge incoming updates into the loaded connection record
	if connection.ConnectionType != "" {
		existing.ConnectionType = connection.ConnectionType
	}
	if connection.RoadType != "" {
		existing.RoadType = connection.RoadType
	}
	if connection.RoadCuttingArea != 0 {
		existing.RoadCuttingArea = connection.RoadCuttingArea
	}
	if connection.ConnectionNo != "" {
		existing.ConnectionNo = connection.ConnectionNo
	}
	if connection.ConnectionHolders != nil {
		existing.ConnectionHolders = connection.ConnectionHolders
	}
	if connection.PlumberInfo != nil {
		existing.PlumberInfo = connection.PlumberInfo
	}
	if connection.RoadCuttingInfo != nil {
		existing.RoadCuttingInfo = connection.RoadCuttingInfo
	}

	// Preserve processInstance action from the incoming request
	existing.ProcessInstance = connection.ProcessInstance

	// Work with the merged connection record
	connection = existing

	// Enrich audit details
	currentTime := time.Now().UnixMilli()
	username := "system"
	if reqInfo.UserInfo != nil {
		username = reqInfo.UserInfo.UserName
	}

	if connection.AuditDetails == nil {
		connection.AuditDetails = &dto.AuditDetails{
			CreatedBy:   username,
			CreatedTime: currentTime,
		}
	}
	connection.AuditDetails.LastModifiedBy = username
	connection.AuditDetails.LastModifiedTime = currentTime

	// If approving or activating, generate connection no if not set
	action := connection.ProcessInstance.Action
	if (action == "APPROVE" || action == "ACTIVATE" || action == "APPROVE_FOR_CONNECTION" || action == "ACTIVATE_CONNECTION") && connection.ConnectionNo == "" {
		var connNo string
		if s.idGenClient != nil {
			connNo, _ = s.idGenClient.GenerateConnectionNo(reqInfo, connection.TenantID)
		}
		if connNo == "" {
			connNo = fmt.Sprintf("SW-CON-%s", uuid.New().String()[:8])
		}
		connection.ConnectionNo = connNo
	}

	// Call workflow
	if s.workflowClient != nil {
		pi, err := s.workflowClient.CallWorkflow(reqInfo, connection)
		if err != nil {
			log.Printf("Warning: Workflow transition failed: %v, using local state machine.", err)
			nextStatus := getLocalNextState(connection.ApplicationStatus, action)
			connection.ApplicationStatus = nextStatus
			connection.Status = nextStatus
		} else if pi != nil {
			if pi.State != nil && pi.State.ApplicationStatus != "" {
				connection.ApplicationStatus = pi.State.ApplicationStatus
				connection.Status = pi.State.ApplicationStatus
			} else {
				connection.ApplicationStatus = pi.Action
				connection.Status = pi.Action
			}
		}
	} else {
		nextStatus := getLocalNextState(connection.ApplicationStatus, action)
		connection.ApplicationStatus = nextStatus
		connection.Status = nextStatus
	}

	// Update DB record
	err = s.repository.UpdateConnection(connection)
	if err != nil {
		return connection, err
	}

	// Publish to Kafka update-sw-connection
	kafkaReq := dto.SewerageConnectionRequest{
		RequestInfo:        reqInfo,
		SewerageConnection: connection,
	}
	_ = s.producer.PublishConnectionUpdated(kafkaReq, appNo)

	return connection, nil
}

func getLocalNextState(currentStatus, action string) string {
	switch action {
	case "SUBMIT_APPLICATION":
		return "PENDING_FOR_DOCUMENT_VERIFICATION"
	case "VERIFY_AND_FORWARD":
		if currentStatus == "PENDING_FOR_DOCUMENT_VERIFICATION" || currentStatus == "INITIATED" || currentStatus == "" {
			return "PENDING_FOR_FIELD_INSPECTION"
		}
		return "PENDING_APPROVAL_FOR_CONNECTION"
	case "APPROVE_FOR_CONNECTION":
		return "PENDING_FOR_PAYMENT"
	case "PAY":
		return "PENDING_FOR_CONNECTION_ACTIVATION"
	case "ACTIVATE_CONNECTION":
		return "CONNECTION_ACTIVATED"
	default:
		return currentStatus
	}
}
