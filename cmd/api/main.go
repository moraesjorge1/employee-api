package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jorgemorais/employee-api/internal/employee"
)

func main() {
	db, err := sql.Open("mysql", "jorge:jorge123@tcp(localhost:3306)/employees_db")
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
	router.HandleFunc("/employees/report", handler.GetReport).Methods("GET") // ← antes do {id}
	router.HandleFunc("/employees/{id}", handler.GetByID).Methods("GET")
	router.HandleFunc("/employees/{id}", handler.Update).Methods("PUT")
	router.HandleFunc("/employees/{id}", handler.Delete).Methods("DELETE")

	fmt.Println("Server running on port 8080!")
	log.Fatal(http.ListenAndServe(":8080", router))
}
