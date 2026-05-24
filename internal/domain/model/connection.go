package model

import "time"

type Connection struct {
	ID                string    `json:"id"`
	TenantID          string    `json:"tenantId"`
	ApplicationNo     string    `json:"applicationNo"`
	ConnectionNo      string    `json:"connectionNo,omitempty"`
	PropertyID        string    `json:"propertyId"`
	ConnectionType    string    `json:"connectionType"`
	ApplicationStatus string    `json:"applicationStatus"`
	Status            string    `json:"status"`
	CreatedBy         string    `json:"createdBy"`
	CreatedTime       time.Time `json:"createdTime"`
	LastModifiedBy    string    `json:"lastModifiedBy"`
	LastModifiedTime  time.Time `json:"lastModifiedTime"`
}
