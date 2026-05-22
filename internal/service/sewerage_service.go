package service

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/BrajK111/sw-services/internal/domain/dto"
	"github.com/BrajK111/sw-services/internal/kafka"
	"github.com/BrajK111/sw-services/internal/repository/postgres"
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
	idGenClient    *IDGenClient
	workflowClient *WorkflowClient
	producer       *kafka.Producer
}

func NewSewerageService(
	repository postgres.SewerageRepository,
	idGen *IDGenClient,
	workflow *WorkflowClient,
	producer *kafka.Producer,
) SewerageService {
	return &sewerageService{
		repository:     repository,
		idGenClient:    idGen,
		workflowClient: workflow,
		producer:       producer,
	}
}

func (s *sewerageService) CreateConnection(
	reqInfo dto.RequestInfo,
	connection dto.SewerageConnection,
) (dto.SewerageConnection, error) {

	id := uuid.New().String()
	connection.ID = id

	// Generate Application Number via IDGen with robust fallback
	var appNo string
	var err error
	if s.idGenClient != nil {
		appNo, err = s.idGenClient.GenerateApplicationNo(reqInfo, connection.TenantID)
	}
	if err != nil || appNo == "" {
		log.Printf("Warning: IDGen unavailable, using fallback UUID for applicationNo: %v", err)
		appNo = fmt.Sprintf("SW-APP-%s", id[:8])
	}

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
			log.Printf("Warning: Workflow transition failed: %v", err)
		} else if pi != nil {
			// update application status if workflow changed it
			connection.ApplicationStatus = pi.Action
			connection.Status = pi.Action
		}
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
	if (action == "APPROVE" || action == "ACTIVATE") && connection.ConnectionNo == "" {
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
			log.Printf("Warning: Workflow transition failed: %v", err)
		} else if pi != nil {
			connection.ApplicationStatus = pi.Action
			connection.Status = pi.Action
		}
	}

	// Update DB record
	err := s.repository.UpdateConnection(connection)
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
