package report

import (
	"time"

	"github.com/jorgemorais/employee-api/internal/employee"
	"go.temporal.io/sdk/workflow"
)

const TaskQueue = "employee-report"

type WorkflowInput struct {
	ReportID string
	Filter   employee.ReportFilter
}

func GenerateReportWorkflow(ctx workflow.Context, input WorkflowInput) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	var act *Activities

	var rep employee.Reporting
	if err := workflow.ExecuteActivity(ctx, act.FetchReportDataActivity, input.Filter).Get(ctx, &rep); err != nil {
		return err
	}

	return workflow.ExecuteActivity(ctx, act.SaveReportActivity, input.ReportID, rep).Get(ctx, nil)
}
