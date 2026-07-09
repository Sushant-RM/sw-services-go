package dto

type SewerageConnectionRequest struct {
	RequestInfo        RequestInfo        `json:"RequestInfo"`
	SewerageConnection SewerageConnection `json:"sewerageConnection"`
	IsCreateCall       bool               `json:"isCreateCall,omitempty"`
	DisconnectRequest  bool               `json:"disconnectRequest,omitempty"`
}

type RequestInfoWrapper struct {
	RequestInfo RequestInfo `json:"RequestInfo"`
}
