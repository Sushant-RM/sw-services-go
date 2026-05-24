package integration

import (
	"net/http"
	"os"
	"time"
)

type PropertyClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewPropertyClient() *PropertyClient {
	host := os.Getenv("EGOV_PROPERTY_SERVICE_HOST")
	if host == "" {
		host = "http://sw-property-services:8080"
	}
	return &PropertyClient{
		baseURL:    host,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}
