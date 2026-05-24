package model

type Property struct {
	PropertyID     string `json:"propertyId"`
	TenantID       string `json:"tenantId"`
	UsageCategory  string `json:"usageCategory"`
	OwnerName      string `json:"ownerName"`
	MobileNumber   string `json:"mobileNumber"`
}

type Billing struct {
	ConsumerCode    string  `json:"consumerCode"`
	TenantID        string  `json:"tenantId"`
	BusinessService string  `json:"businessService"`
	TaxAmount       float64 `json:"taxAmount"`
}

type Payment struct {
	PaymentID     string  `json:"paymentId"`
	TenantID      string  `json:"tenantId"`
	TransactionID string  `json:"transactionId"`
	AmountPaid    float64 `json:"amountPaid"`
	Status        string  `json:"status"`
}

type User struct {
	UUID         string `json:"uuid"`
	UserName     string `json:"userName"`
	Name         string `json:"name"`
	MobileNumber string `json:"mobileNumber"`
	Active       bool   `json:"active"`
	TenantID     string `json:"tenantId"`
}
