package workflow

type ProcessInstance struct {
	BusinessID string `json:"businessId"`

	BusinessService string `json:"businessService"`

	Action string `json:"action"`

	ModuleName string `json:"moduleName"`

	State string `json:"state"`
}
