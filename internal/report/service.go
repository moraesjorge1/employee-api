package report

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/jorgemorais/employee-api/internal/employee"
	"go.temporal.io/sdk/client"
)

type Service struct {
	client client.Client
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

func newUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
