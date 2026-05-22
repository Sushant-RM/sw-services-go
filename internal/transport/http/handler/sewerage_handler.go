package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/BrajK111/sw-services/internal/domain/dto"
	"github.com/BrajK111/sw-services/internal/service"
	"github.com/BrajK111/sw-services/internal/transport/http/response"
)

type SewerageHandler struct {
	service service.SewerageService
}

func NewSewerageHandler(
	service service.SewerageService,
) *SewerageHandler {
	return &SewerageHandler{
		service: service,
	}
}

// CreateConnection godoc
// @Summary Create sewerage connection
// @Description Creates a new sewerage connection
// @Tags Sewerage
// @Accept json
// @Produce json
// @Success 200 {object} dto.SewerageConnectionResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /sw-services/swc/_create [post]
func (h *SewerageHandler) CreateConnection(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request dto.SewerageConnectionRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.WriteError(
			w,
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request payload: "+err.Error(),
		)
		return
	}

	createdConn, err := h.service.CreateConnection(
		request.RequestInfo,
		request.SewerageConnection,
	)
	if err != nil {
		response.WriteError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			err.Error(),
		)
		return
	}

	createResponse := dto.SewerageConnectionResponse{
		ResponseInfo: dto.ResponseInfo{
			APIID:    request.RequestInfo.APIID,
			Ver:      request.RequestInfo.Ver,
			Ts:       time.Now().UnixMilli(),
			ResMsgID: "uief87324",
			MsgID:    request.RequestInfo.MsgID,
			Status:   "successful",
		},
		SewerageConnections: []dto.SewerageConnection{
			createdConn,
		},
		TotalCount: 1,
	}

	response.WriteJSON(
		w,
		http.StatusOK,
		createResponse,
	)
}

// SearchConnection godoc
// @Summary Search sewerage connections
// @Description Searches sewerage connections
// @Tags Sewerage
// @Produce json
// @Param tenantId query string false "Tenant ID"
// @Param applicationNumber query string false "Application Number"
// @Success 200 {object} dto.SewerageConnectionResponse
// @Failure 500 {object} map[string]interface{}
// @Router /sw-services/swc/_search [post]
func (h *SewerageHandler) SearchConnection(
	w http.ResponseWriter,
	r *http.Request,
) {
	tenantID := r.URL.Query().Get("tenantId")
	applicationNumber := r.URL.Query().Get("applicationNumber")
	if applicationNumber == "" {
		// Fallback to applicationNo query parameter if passed
		applicationNumber = r.URL.Query().Get("applicationNo")
	}

	var reqInfo dto.RequestInfo

	// Read RequestInfo from body if request is POST
	if r.Method == http.MethodPost {
		var wrapper dto.RequestInfoWrapper
		_ = json.NewDecoder(r.Body).Decode(&wrapper)
		reqInfo = wrapper.RequestInfo
	}

	connections, err := h.service.SearchConnection(
		reqInfo,
		tenantID,
		applicationNumber,
	)
	if err != nil {
		response.WriteError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			err.Error(),
		)
		return
	}

	searchResponse := dto.SewerageConnectionResponse{
		ResponseInfo: dto.ResponseInfo{
			APIID:    reqInfo.APIID,
			Ver:      reqInfo.Ver,
			Ts:       time.Now().UnixMilli(),
			ResMsgID: "uief87324",
			MsgID:    reqInfo.MsgID,
			Status:   "successful",
		},
		SewerageConnections: connections,
		TotalCount:          len(connections),
	}

	response.WriteJSON(
		w,
		http.StatusOK,
		searchResponse,
	)
}

// UpdateConnection godoc
// @Summary Update sewerage connection
// @Description Updates an existing sewerage connection
// @Tags Sewerage
// @Accept json
// @Produce json
// @Success 200 {object} dto.SewerageConnectionResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /sw-services/swc/_update [post]
func (h *SewerageHandler) UpdateConnection(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request dto.SewerageConnectionRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.WriteError(
			w,
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request payload: "+err.Error(),
		)
		return
	}

	updatedConn, err := h.service.UpdateConnection(
		request.RequestInfo,
		request.SewerageConnection,
	)
	if err != nil {
		response.WriteError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			err.Error(),
		)
		return
	}

	updateResponse := dto.SewerageConnectionResponse{
		ResponseInfo: dto.ResponseInfo{
			APIID:    request.RequestInfo.APIID,
			Ver:      request.RequestInfo.Ver,
			Ts:       time.Now().UnixMilli(),
			ResMsgID: "uief87324",
			MsgID:    request.RequestInfo.MsgID,
			Status:   "successful",
		},
		SewerageConnections: []dto.SewerageConnection{
			updatedConn,
		},
		TotalCount: 1,
	}

	response.WriteJSON(
		w,
		http.StatusOK,
		updateResponse,
	)
}

// CountConnections godoc
// @Summary Count sewerage connections
// @Description Counts sewerage connections matching criteria
// @Tags Sewerage
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /sw-services/swc/_count [post]
func (h *SewerageHandler) CountConnections(
	w http.ResponseWriter,
	r *http.Request,
) {
	tenantID := r.URL.Query().Get("tenantId")
	applicationNumber := r.URL.Query().Get("applicationNumber")

	var reqInfo dto.RequestInfo
	if r.Method == http.MethodPost {
		var wrapper dto.RequestInfoWrapper
		_ = json.NewDecoder(r.Body).Decode(&wrapper)
		reqInfo = wrapper.RequestInfo
	}

	connections, err := h.service.SearchConnection(
		reqInfo,
		tenantID,
		applicationNumber,
	)
	if err != nil {
		response.WriteError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			err.Error(),
		)
		return
	}

	countResponse := map[string]interface{}{
		"ResponseInfo": dto.ResponseInfo{
			APIID:    reqInfo.APIID,
			Ver:      reqInfo.Ver,
			Ts:       time.Now().UnixMilli(),
			ResMsgID: "uief87324",
			MsgID:    reqInfo.MsgID,
			Status:   "successful",
		},
		"count": len(connections),
	}

	response.WriteJSON(
		w,
		http.StatusOK,
		countResponse,
	)
}

// EncryptOldData godoc
// @Summary Encrypt old connection data
// @Description Throws error stating privacy/encryption is disabled
// @Tags Sewerage
// @Produce json
// @Success 400 {object} map[string]interface{}
// @Router /sw-services/swc/_encryptOldData [post]
func (h *SewerageHandler) EncryptOldData(
	w http.ResponseWriter,
	r *http.Request,
) {
	response.WriteError(
		w,
		http.StatusBadRequest,
		"EG_SW_ENC_OLD_DATA_ERROR",
		"Privacy disabled: The encryption of old data is disabled",
	)
}
