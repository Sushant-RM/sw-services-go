package dto

type SewerageConnectionResponse struct {
	ResponseInfo        ResponseInfo         `json:"ResponseInfo"`
	SewerageConnections []SewerageConnection `json:"SewerageConnections"`
	TotalCount          int                  `json:"totalCount,omitempty"`
}

type CreateSewerageConnectionResponse struct {
	ResponseInfo        ResponseInfo         `json:"ResponseInfo"`
	SewerageConnections []SewerageConnection `json:"SewerageConnections"`
}

type SearchSewerageConnectionResponse struct {
	ResponseInfo        ResponseInfo         `json:"ResponseInfo"`
	SewerageConnections []SewerageConnection `json:"SewerageConnections"`
	TotalCount          int                  `json:"TotalCount"`
}
