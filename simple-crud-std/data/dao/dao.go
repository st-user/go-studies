package dao

import (
	"database/sql"
	"fmt"
	"simple-crud-std/data/entity"
)

type EmployeeDAO struct {
	db *sql.DB
}

func NewEmployeeDAO(db *sql.DB) EmployeeDAO {
	return EmployeeDAO{db}
}

func (e EmployeeDAO) FindAll() ([]entity.Employee, error) {
	employees := []entity.Employee{}
	rows, err := e.db.Query("SELECT * FROM employees")
	if err != nil {
		return nil, fmt.Errorf("error at query@'FindAll': %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var employee entity.Employee
		if err := rows.Scan(
			&employee.ID,
			&employee.Name,
			&employee.StartDate,
		); err != nil {
			return nil, fmt.Errorf("error at scan@'FindAll': %w", err)
		}
		employees = append(employees, employee)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error at err@'FindAll': %w", err)
	}

	return employees, nil
}

func (e EmployeeDAO) FindByID(id int64) (entity.Employee, error) {

	var employee entity.Employee

	row := e.db.QueryRow("SELECT * FROM employees WHERE id = ?", id)
	if err := row.Scan(
		&employee.ID,
		&employee.Name,
		&employee.StartDate,
	); err != nil {
		if err == sql.ErrNoRows {
			return employee, entity.NewNotFoundError(
				fmt.Sprintf("employee not found for: %d", id),
			)
		}
		return employee, fmt.Errorf("error at scan@'FindByID': %w", err)
	}

	return employee, nil
}

func (e EmployeeDAO) Save(employee entity.Employee) (entity.Employee, error) {
	result, err := e.db.Exec(
		"INSERT INTO employees (name, start_date) VALUES (?, ?)",
		employee.Name,
		employee.StartDate)

	if err != nil {
		return entity.Employee{}, fmt.Errorf("error at exec@'Save': %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return entity.Employee{}, fmt.Errorf("error at last_inserted_id@'Save': %w", err)
	}

	return entity.Employee{
		ID:        id,
		Name:      employee.Name,
		StartDate: employee.StartDate,
	}, nil
}

func (e EmployeeDAO) Update(employee entity.Employee) error {
	result, err := e.db.Exec(
		"UPDATE employees SET name = ?, start_date = ? WHERE id = ?",
		employee.Name,
		employee.StartDate,
		employee.ID,
	)

	cnt, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error at exec@'Update': %w", err)
	}

	if cnt == 0 {
		return entity.NewNotFoundError(
			fmt.Sprintf("employee not found for: %d", employee.ID),
		)
	}

	return nil
}

func (e EmployeeDAO) DeleteByID(id int64) error {
	_, err := e.db.Exec(
		"DELETE FROM employees WHERE id = ?",
		id,
	)

	if err != nil {
		return fmt.Errorf("error at exec@'DeleteByID': %w", err)
	}

	return nil
}
