package integration

import (
	"net/http"
	"os"
	"time"
)

type BillingClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewBillingClient() *BillingClient {
	host := os.Getenv("EGOV_BILLING_SERVICE_HOST")
	if host == "" {
		host = "http://sw-billing-service:8080"
	}
	return &BillingClient{
		baseURL:    host,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}
