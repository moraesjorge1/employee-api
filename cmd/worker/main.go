package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jorgemorais/employee-api/internal/employee"
	"github.com/jorgemorais/employee-api/internal/report"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN environment variable is required")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatal("unable to create Temporal client: ", err)
	}
	defer c.Close()

	repo := employee.NewRepository(db)
	svc := employee.NewService(repo)
	activities := report.NewActivities(svc)

	w := worker.New(c, report.TaskQueue, worker.Options{})
	w.RegisterWorkflow(report.GenerateReportWorkflow)
	w.RegisterActivity(activities)

	log.Println("Worker started on task queue:", report.TaskQueue)
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatal("worker error: ", err)
	}
}
