package service

type WorkflowService interface{}

type workflowService struct{}

func NewWorkflowService() WorkflowService {
	return &workflowService{}
}
