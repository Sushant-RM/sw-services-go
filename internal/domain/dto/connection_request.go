package dto

type CreateSewerageConnectionRequest struct {
	RequestInfo                RequestInfo        `json:"RequestInfo"`
	SewerageConnection         SewerageConnection `json:"SewerageConnection"`
	IsCreateCall               bool               `json:"isCreateCall,omitempty"`
	IsOldDataEncryptionRequest bool               `json:"isOldDataEncryptionRequest,omitempty"`
	DisconnectRequest          bool               `json:"disconnectRequest,omitempty"`
}

type SewerageConnectionRequest struct {
	RequestInfo                RequestInfo        `json:"RequestInfo"`
	SewerageConnection         SewerageConnection `json:"SewerageConnection"`
	IsCreateCall               bool               `json:"isCreateCall,omitempty"`
	IsOldDataEncryptionRequest bool               `json:"isOldDataEncryptionRequest,omitempty"`
	DisconnectRequest          bool               `json:"disconnectRequest,omitempty"`
}

type SearchSewerageConnectionRequest struct {
	RequestInfo RequestInfo `json:"RequestInfo"`
}

type RequestInfoWrapper struct {
	RequestInfo RequestInfo `json:"RequestInfo"`
}
