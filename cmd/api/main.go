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

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to database")

	// 3. cria a tabela se não existir
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

	repo := employee.NewRepository(db)
	service := employee.NewService(repo)
	handler := employee.NewHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/employees", handler.Create).Methods("POST")
	router.HandleFunc("/employees", handler.GetAll).Methods("GET")
	router.HandleFunc("/employees/report", handler.GetReport).Methods("GET")
	router.HandleFunc("/employees/{id}", handler.GetByID).Methods("GET")
	router.HandleFunc("/employees/{id}", handler.Update).Methods("PUT")
	router.HandleFunc("/employees/{id}", handler.Delete).Methods("DELETE")

	fmt.Println("Server running on port 8080!")
	log.Fatal(http.ListenAndServe(":8080", router))
}
