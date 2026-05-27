package employee

import (
	"database/sql"
	"errors"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(emp Employee) (Employee, error) {
	result, err := r.db.Exec(
		"INSERT INTO employees (name, position, salary, type) VALUES (?, ?, ?, ?)",
		emp.Name, emp.Position, emp.Salary, emp.Type,
	)
	if err != nil {
		return Employee{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Employee{}, err
	}

	emp.ID = id
	return emp, nil
}

func (r *Repository) GetAll() ([]Employee, error) {
	rows, err := r.db.Query("SELECT id, name, position, salary, type FROM employees")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := []Employee{}
	for rows.Next() {
		var emp Employee
		err := rows.Scan(&emp.ID, &emp.Name, &emp.Position, &emp.Salary, &emp.Type)
		if err != nil {
			return nil, err
		}
		employees = append(employees, emp)
	}

	return employees, nil
}

func (r *Repository) GetByID(id int64) (Employee, error) {
	var emp Employee
	err := r.db.QueryRow(
		"SELECT id, name, position, salary, type FROM employees WHERE id = ?", id,
	).Scan(&emp.ID, &emp.Name, &emp.Position, &emp.Salary, &emp.Type)
	if errors.Is(err, sql.ErrNoRows) {
		return Employee{}, ErrNotFound
	}
	if err != nil {
		return Employee{}, err
	}

	return emp, nil
}

func (r *Repository) Update(emp Employee) (Employee, error) {
	result, err := r.db.Exec(
		"UPDATE employees SET name = ?, position = ?, salary = ?, type = ? WHERE id = ?",
		emp.Name, emp.Position, emp.Salary, emp.Type, emp.ID,
	)
	if err != nil {
		return Employee{}, err
	}

	// verifica se alguma linha foi afetada
	rows, err := result.RowsAffected()
	if err != nil {
		return Employee{}, err
	}
	if rows == 0 {
		return Employee{}, fmt.Errorf("employee with id %d not found", emp.ID)
	}

	return emp, nil
}

func (r *Repository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM employees WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("employee with id %d not found", id)
	}

	return nil
}
