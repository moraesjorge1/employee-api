package report

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/jorgemorais/employee-api/internal/employee"
)

const reportsDir = "reports"

type Activities struct {
	svc *employee.Service
}

func NewActivities(svc *employee.Service) *Activities {
	return &Activities{svc: svc}
}

func (a *Activities) FetchReportDataActivity(ctx context.Context, filter employee.ReportFilter) (employee.Reporting, error) {
	return a.svc.GenerateReport(filter)
}

func (a *Activities) SaveReportActivity(ctx context.Context, reportID string, rep employee.Reporting) error {
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		return fmt.Errorf("creating reports dir: %w", err)
	}
	data, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fmt.Sprintf("%s/%s.json", reportsDir, reportID), data, 0644)
}

func readReport(reportID string) (*employee.Reporting, bool, error) {
	data, err := os.ReadFile(fmt.Sprintf("%s/%s.json", reportsDir, reportID))
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	var rep employee.Reporting
	if err := json.Unmarshal(data, &rep); err != nil {
		return nil, false, err
	}
	return &rep, true, nil
}
