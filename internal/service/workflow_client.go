// WORKFLOW CLIENT: Drop into internal/service/workflow_client.go
// Calls DIGIT's egov-workflow-v2 service for state transitions

package service

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

type WorkflowClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewWorkflowClient() *WorkflowClient {
	host := os.Getenv("WORKFLOW_HOST")
	if host == "" {
		host = "http://egov-workflow-v2:8080"
	}
	return &WorkflowClient{
		baseURL:    host,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// CallWorkflow transitions workflow state - call this on EVERY create/update
func (c *WorkflowClient) CallWorkflow(reqInfo dto.RequestInfo, conn dto.SewerageConnection) (*dto.ProcessInstance, error) {
	if conn.ProcessInstance.Action == "" {
		return nil, nil // workflow not triggered if no ProcessInstance in request
	}

	pi := conn.ProcessInstance
	appNo := conn.ApplicationNo
	if appNo == "" {
		appNo = conn.ApplicationNumber
	}
	pi.BusinessID = appNo
	pi.TenantID = conn.TenantID
	pi.ModuleName = "SW"
	if pi.BusinessService == "" || pi.BusinessService == "SW" {
		if conn.ApplicationType == "MODIFY_CONNECTION" {
			pi.BusinessService = "ModifySWConnection"
		} else if conn.ApplicationType == "DISCONNECT_CONNECTION" {
			pi.BusinessService = "DisconnectSWConnection"
		} else {
			pi.BusinessService = "NewSW1"
		}
	}

	payload := dto.WorkflowRequest{
		RequestInfo:      reqInfo,
		ProcessInstances: []dto.ProcessInstance{pi},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal workflow request: %w", err)
	}

	resp, err := c.httpClient.Post(
		c.baseURL+"/egov-workflow-v2/egov-wf/process/_transition",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("workflow http call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("workflow returned %d: %s", resp.StatusCode, string(b))
	}

	var result dto.WorkflowResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode workflow response: %w", err)
	}

	if len(result.ProcessInstances) == 0 {
		return nil, fmt.Errorf("workflow returned empty response")
	}

	return &result.ProcessInstances[0], nil
}

// GetWorkflowStatus retrieves current workflow status for an application
func (c *WorkflowClient) GetWorkflowStatus(reqInfo dto.RequestInfo, tenantId, businessId string) (string, error) {
	url := fmt.Sprintf("%s/egov-workflow-v2/egov-wf/process/_search?tenantId=%s&businessIds=%s",
		c.baseURL, tenantId, businessId)

	reqBody := struct {
		RequestInfo dto.RequestInfo `json:"RequestInfo"`
	}{RequestInfo: reqInfo}

	body, _ := json.Marshal(reqBody)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	// Extract status from response - adjust based on actual workflow response structure
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return "", err
	}

	if pis, ok := raw["ProcessInstances"].([]interface{}); ok && len(pis) > 0 {
		if pi, ok := pis[0].(map[string]interface{}); ok {
			if state, ok := pi["state"].(map[string]interface{}); ok {
				return fmt.Sprintf("%v", state["applicationStatus"]), nil
			}
		}
	}
	return "UNKNOWN", nil
}
