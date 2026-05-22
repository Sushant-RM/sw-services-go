package workflow

import (
	"log"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) StartWorkflow(
	process ProcessInstance,
) error {

	log.Printf(
		"Workflow started | BusinessID=%s | Action=%s | State=%s",
		process.BusinessID,
		process.Action,
		process.State,
	)

	return nil
}
