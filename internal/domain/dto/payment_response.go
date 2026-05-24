package dto

type PaymentResponse struct {
	ResponseInfo ResponseInfo `json:"ResponseInfo"`
	Payments     []interface{} `json:"Payments"`
}
