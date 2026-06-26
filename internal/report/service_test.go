package report

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jorgemorais/employee-api/internal/employee"
	enumspb "go.temporal.io/api/enums/v1"
	workflowpb "go.temporal.io/api/workflow/v1"
	workflowservice "go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

// --- mock do Temporal client ---

type mockTemporalClient struct {
	executeWorkflowFn       func(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error)
	describeWorkflowFn      func(ctx context.Context, workflowID, runID string) (*workflowservice.DescribeWorkflowExecutionResponse, error)
}

func (m *mockTemporalClient) ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error) {
	return m.executeWorkflowFn(ctx, options, workflow, args...)
}

func (m *mockTemporalClient) DescribeWorkflowExecution(ctx context.Context, workflowID, runID string) (*workflowservice.DescribeWorkflowExecutionResponse, error) {
	return m.describeWorkflowFn(ctx, workflowID, runID)
}

// mock do WorkflowRun retornado pelo ExecuteWorkflow
type mockWorkflowRun struct{}

func (m *mockWorkflowRun) GetID() string     { return "mock-id" }
func (m *mockWorkflowRun) GetRunID() string  { return "mock-run-id" }
func (m *mockWorkflowRun) Get(ctx context.Context, valuePtr interface{}) error { return nil }
func (m *mockWorkflowRun) GetWithOptions(ctx context.Context, valuePtr interface{}, opts client.WorkflowRunGetOptions) error {
	return nil
}

// helper: cria Service com mock e diretório temporário
func newTestService(t *testing.T, c temporalClient) *Service {
	t.Helper()
	dir := t.TempDir()
	reportsDir = dir
	t.Cleanup(func() { reportsDir = "reports" })
	return &Service{client: c}
}

// helper: salva um report em disco no diretório de teste
func writeTestReport(t *testing.T, reportID string, rep employee.Reporting) {
	t.Helper()
	data, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(reportsDir, fmt.Sprintf("%s.json", reportID))
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
}

// --- DownloadReport ---

func TestDownloadReport_FileExists(t *testing.T) {
	svc := newTestService(t, nil)
	rep := employee.Reporting{TotalEmployees: 1, TotalSalary: 5000, AverageSalary: 5000}
	writeTestReport(t, "abc-123", rep)

	data, err := svc.DownloadReport("abc-123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty data")
	}
}

func TestDownloadReport_FileNotFound(t *testing.T) {
	svc := newTestService(t, nil)

	_, err := svc.DownloadReport("nao-existe")
	if !errors.Is(err, ErrReportNotFound) {
		t.Fatalf("expected ErrReportNotFound, got %v", err)
	}
}

// --- GetStatus ---

func TestGetStatus_FileExists_ReturnsDone(t *testing.T) {
	svc := newTestService(t, nil)
	writeTestReport(t, "done-id", employee.Reporting{TotalEmployees: 2})

	status, err := svc.GetStatus(context.Background(), "done-id")
	if err != nil {
		t.Fatal(err)
	}
	if status != StatusDone {
		t.Fatalf("expected done, got %s", status)
	}
}

func TestGetStatus_WorkflowRunning_ReturnsPending(t *testing.T) {
	mock := &mockTemporalClient{
		describeWorkflowFn: func(ctx context.Context, workflowID, runID string) (*workflowservice.DescribeWorkflowExecutionResponse, error) {
			return &workflowservice.DescribeWorkflowExecutionResponse{
				WorkflowExecutionInfo: &workflowpb.WorkflowExecutionInfo{
					Status: enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING,
				},
			}, nil
		},
	}
	svc := newTestService(t, mock)

	status, err := svc.GetStatus(context.Background(), "running-id")
	if err != nil {
		t.Fatal(err)
	}
	if status != StatusPending {
		t.Fatalf("expected pending, got %s", status)
	}
}

func TestGetStatus_WorkflowFailed_ReturnsFailed(t *testing.T) {
	mock := &mockTemporalClient{
		describeWorkflowFn: func(ctx context.Context, workflowID, runID string) (*workflowservice.DescribeWorkflowExecutionResponse, error) {
			return &workflowservice.DescribeWorkflowExecutionResponse{
				WorkflowExecutionInfo: &workflowpb.WorkflowExecutionInfo{
					Status: enumspb.WORKFLOW_EXECUTION_STATUS_FAILED,
				},
			}, nil
		},
	}
	svc := newTestService(t, mock)

	status, err := svc.GetStatus(context.Background(), "failed-id")
	if err != nil {
		t.Fatal(err)
	}
	if status != StatusFailed {
		t.Fatalf("expected failed, got %s", status)
	}
}

func TestGetStatus_WorkflowCompletedNoFile_ReturnsFailed(t *testing.T) {
	mock := &mockTemporalClient{
		describeWorkflowFn: func(ctx context.Context, workflowID, runID string) (*workflowservice.DescribeWorkflowExecutionResponse, error) {
			return &workflowservice.DescribeWorkflowExecutionResponse{
				WorkflowExecutionInfo: &workflowpb.WorkflowExecutionInfo{
					Status: enumspb.WORKFLOW_EXECUTION_STATUS_COMPLETED,
				},
			}, nil
		},
	}
	svc := newTestService(t, mock)

	// workflow completou mas SaveReportActivity falhou — arquivo não existe
	status, err := svc.GetStatus(context.Background(), "completed-no-file")
	if err != nil {
		t.Fatal(err)
	}
	if status != StatusFailed {
		t.Fatalf("expected failed, got %s", status)
	}
}

func TestGetStatus_TemporalError_ReturnsFailed(t *testing.T) {
	mock := &mockTemporalClient{
		describeWorkflowFn: func(ctx context.Context, workflowID, runID string) (*workflowservice.DescribeWorkflowExecutionResponse, error) {
			return nil, errors.New("temporal unavailable")
		},
	}
	svc := newTestService(t, mock)

	_, err := svc.GetStatus(context.Background(), "some-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- StartReport ---

func TestStartReport_TriggersWorkflow(t *testing.T) {
	var capturedOptions client.StartWorkflowOptions
	var capturedInput WorkflowInput

	mock := &mockTemporalClient{
		executeWorkflowFn: func(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error) {
			capturedOptions = options
			if len(args) > 0 {
				if input, ok := args[0].(WorkflowInput); ok {
					capturedInput = input
				}
			}
			return &mockWorkflowRun{}, nil
		},
	}
	svc := newTestService(t, mock)

	filter := employee.ReportFilter{Type: "fulltime", MinSalary: 5000}
	reportID, err := svc.StartReport(filter)
	if err != nil {
		t.Fatal(err)
	}
	if reportID == "" {
		t.Fatal("expected non-empty report ID")
	}
	if capturedOptions.TaskQueue != TaskQueue {
		t.Fatalf("expected task queue %s, got %s", TaskQueue, capturedOptions.TaskQueue)
	}
	if capturedInput.ReportID != reportID {
		t.Fatalf("expected workflow input ReportID %s, got %s", reportID, capturedInput.ReportID)
	}
	if capturedInput.Filter.Type != "fulltime" {
		t.Fatalf("expected filter type fulltime, got %s", capturedInput.Filter.Type)
	}
}

func TestStartReport_TemporalError(t *testing.T) {
	mock := &mockTemporalClient{
		executeWorkflowFn: func(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error) {
			return nil, errors.New("temporal unavailable")
		},
	}
	svc := newTestService(t, mock)

	_, err := svc.StartReport(employee.ReportFilter{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
