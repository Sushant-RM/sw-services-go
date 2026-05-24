package model

type Workflow struct {
	ProcessInstanceID string `json:"processInstanceId"`
	BusinessID        string `json:"businessId"`
	Action            string `json:"action"`
	Status            string `json:"status"`
	Comment           string `json:"comment"`
}
