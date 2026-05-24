package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/BrajK111/sw-services/internal/domain/dto"
)

type IDGenClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewIDGenClient() *IDGenClient {
	host := os.Getenv("EGOV_IDGEN_HOST")
	if host == "" {
		host = "http://egov-idgen:8080"
	}
	return &IDGenClient{
		baseURL:    host,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *IDGenClient) GenerateID(reqInfo dto.RequestInfo, idName, format, tenantId string) (string, error) {
	payload := dto.IDGenRequest{
		RequestInfo: reqInfo,
		IDRequests: []dto.IDRequest{
			{
				IDName:   idName,
				Format:   format,
				TenantID: tenantId,
				Count:    1,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal idgen request: %w", err)
	}

	resp, err := c.httpClient.Post(
		c.baseURL+"/egov-idgen/id/_generate",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", fmt.Errorf("idgen http call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("idgen returned %d: %s", resp.StatusCode, string(b))
	}

	var result dto.IDGenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode idgen response: %w", err)
	}

	if len(result.IDResponses) == 0 {
		return "", fmt.Errorf("idgen returned empty response")
	}

	return result.IDResponses[0].ID, nil
}

func (c *IDGenClient) GenerateApplicationNo(reqInfo dto.RequestInfo, tenantId string) (string, error) {
	return c.GenerateID(
		reqInfo,
		"sewerageservice.application.id",
		"SW_AP/[CITY.CODE]/[fy:yyyy-yy]/[SEQ_SW_APP_[TENANT_ID]]",
		tenantId,
	)
}

func (c *IDGenClient) GenerateConnectionNo(reqInfo dto.RequestInfo, tenantId string) (string, error) {
	return c.GenerateID(
		reqInfo,
		"sewerageservice.connection.id",
		"SW/[CITY.CODE]/[fy:yyyy-yy]/[SEQ_SW_CON_[TENANT_ID]]",
		tenantId,
	)
}
