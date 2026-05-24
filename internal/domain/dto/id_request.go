package dto

type IDGenRequest struct {
	RequestInfo RequestInfo `json:"RequestInfo"`
	IDRequests  []IDRequest `json:"idRequests"`
}

type IDRequest struct {
	IDName   string `json:"idName"`
	Format   string `json:"format"`
	TenantID string `json:"tenantId"`
	Count    int    `json:"count"`
}
