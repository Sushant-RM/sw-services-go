package dto

type PaymentRequest struct {
	RequestInfo RequestInfo `json:"RequestInfo"`
	Payment     interface{} `json:"Payment"`
}
