package report

import "errors"

var ErrReportNotFound = errors.New("report not found")

type Status string

const (
	StatusPending Status = "pending"
	StatusDone    Status = "done"
	StatusFailed  Status = "failed"
)

type StatusResponse struct {
	ReportID string `json:"report_id"`
	Status   Status `json:"status"`
}
