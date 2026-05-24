package integration

import (
	"net/http"
	"os"
	"time"
)

type UserClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewUserClient() *UserClient {
	host := os.Getenv("EGOV_USER_HOST")
	if host == "" {
		host = "http://sw-egov-user:8080"
	}
	return &UserClient{
		baseURL:    host,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}
