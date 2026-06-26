package employee

import (
	"database/sql"
	"errors"
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

func (r *Repository) GetFiltered(f ReportFilter) ([]Employee, error) {
	query := "SELECT id, name, position, salary, type FROM employees WHERE 1=1"
	args := []interface{}{}

	if f.Type != "" {
		query += " AND type = ?"
		args = append(args, f.Type)
	}
	if f.Position != "" {
		query += " AND position = ?"
		args = append(args, f.Position)
	}
	if f.MinSalary > 0 {
		query += " AND salary >= ?"
		args = append(args, f.MinSalary)
	}
	if f.MaxSalary > 0 {
		query += " AND salary <= ?"
		args = append(args, f.MaxSalary)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := []Employee{}
	for rows.Next() {
		var emp Employee
		if err := rows.Scan(&emp.ID, &emp.Name, &emp.Position, &emp.Salary, &emp.Type); err != nil {
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
		return Employee{}, ErrNotFound
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
		return ErrNotFound
	}

	return nil
}
