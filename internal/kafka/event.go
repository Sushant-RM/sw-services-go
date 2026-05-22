package kafka

type SewerageConnectionCreatedEvent struct {
	ApplicationNumber string `json:"applicationNumber"`

	TenantID string `json:"tenantId"`

	Status string `json:"status"`

	EventType string `json:"eventType"`
}
