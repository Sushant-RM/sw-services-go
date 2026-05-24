package consumer

type SewerageConnectionCreatedEvent struct {
	ApplicationNumber string `json:"applicationNumber"`
	TenantID          string `json:"tenantId"`
	Status            string `json:"status"`
	EventType         string `json:"eventType"`
}

// Consumer placeholder for target structure alignment
type Consumer struct{}
