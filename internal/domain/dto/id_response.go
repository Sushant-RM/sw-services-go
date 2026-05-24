package dto

type IDGenResponse struct {
	ResponseInfo ResponseInfo `json:"ResponseInfo"`
	IDResponses  []IDResponse `json:"idResponses"`
}

type IDResponse struct {
	IDName string `json:"idName"`
	Format string `json:"format"`
	ID     string `json:"id"`
}
