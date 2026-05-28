package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jorgemorais/employee-api/internal/employee"
	"github.com/jorgemorais/employee-api/internal/report"
	"go.temporal.io/sdk/client"
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
	fmt.Println("Connected to database")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS employees (
			id       INT AUTO_INCREMENT PRIMARY KEY,
			name     VARCHAR(100) NOT NULL,
			position VARCHAR(100) NOT NULL,
			salary   DECIMAL(10,2) NOT NULL,
			type     VARCHAR(20) NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table ready!")

	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatal("unable to create Temporal client: ", err)
	}
	defer temporalClient.Close()

	repo := employee.NewRepository(db)
	svc := employee.NewService(repo)
	empHandler := employee.NewHandler(svc)

	reportSvc := report.NewService(temporalClient)
	reportHandler := report.NewHandler(reportSvc)

	router := mux.NewRouter()
	router.HandleFunc("/employees", empHandler.Create).Methods("POST")
	router.HandleFunc("/employees", empHandler.GetAll).Methods("GET")
	router.HandleFunc("/employees/{id}", empHandler.GetByID).Methods("GET")
	router.HandleFunc("/employees/{id}", empHandler.Update).Methods("PUT")
	router.HandleFunc("/employees/{id}", empHandler.Delete).Methods("DELETE")

	router.HandleFunc("/employees/report", reportHandler.StartReport).Methods("POST")
	router.HandleFunc("/employees/report/{report_id}", reportHandler.GetReport).Methods("GET")

	fmt.Println("Server running on port 8080!")
	log.Fatal(http.ListenAndServe(":8080", router))
}
