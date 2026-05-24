package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func (c *WorkflowClient) CallWorkflow(reqInfo dto.RequestInfo, conn dto.SewerageConnection) (*dto.ProcessInstance, error) {
	if conn.ProcessInstance.Action == "" {
		return nil, nil
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

	debugBytes, _ := json.Marshal(result)
	log.Printf("[DEBUG] Workflow transition response: %s", string(debugBytes))

	if len(result.ProcessInstances) == 0 {
		return nil, fmt.Errorf("workflow returned empty response")
	}

	return &result.ProcessInstances[0], nil
}

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
