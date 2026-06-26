package report

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"os"

	"github.com/jorgemorais/employee-api/internal/employee"
	enumspb "go.temporal.io/api/enums/v1"
	workflowservice "go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

type temporalClient interface {
	ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error)
	DescribeWorkflowExecution(ctx context.Context, workflowID, runID string) (*workflowservice.DescribeWorkflowExecutionResponse, error)
}

type Service struct {
	client temporalClient
}

func NewService(c client.Client) *Service {
	return &Service{client: c}
}

func (s *Service) StartReport(filter employee.ReportFilter) (string, error) {
	reportID := newUUID()
	_, err := s.client.ExecuteWorkflow(
		context.Background(),
		client.StartWorkflowOptions{
			ID:        "report-" + reportID,
			TaskQueue: TaskQueue,
		},
		GenerateReportWorkflow,
		WorkflowInput{ReportID: reportID, Filter: filter},
	)
	if err != nil {
		return "", err
	}
	return reportID, nil
}

func (s *Service) GetReport(reportID string) (*employee.Reporting, bool, error) {
	return readReport(reportID)
}

func (s *Service) GetStatus(ctx context.Context, reportID string) (Status, error) {
	_, found, err := readReport(reportID)
	if err != nil {
		return StatusFailed, err
	}
	if found {
		return StatusDone, nil
	}

	resp, err := s.client.DescribeWorkflowExecution(ctx, "report-"+reportID, "")
	if err != nil {
		return StatusFailed, err
	}

	switch resp.WorkflowExecutionInfo.Status {
	case enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING:
		return StatusPending, nil
	case enumspb.WORKFLOW_EXECUTION_STATUS_COMPLETED:
		// workflow completou mas arquivo não existe: SaveReportActivity falhou
		return StatusFailed, nil
	default:
		return StatusFailed, nil
	}
}

func (s *Service) DownloadReport(reportID string) ([]byte, error) {
	data, err := os.ReadFile(fmt.Sprintf("%s/%s.json", reportsDir, reportID))
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrReportNotFound
	}
	return data, err
}

func newUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
