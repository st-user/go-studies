package employee

import (
	"context"
	"errors"
	entities "example/simple-crud/pkg/entities/employee"
	"example/simple-crud/pkg/interfaces/gateways/db/internal"
	usecases "example/simple-crud/pkg/usecases/employee"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) FindByID(ctx context.Context, id int64) (usecases.OneRepositoryResult, error) {
	var employee entities.Employee
	res := r.db.First(&employee, id)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return usecases.OneRepositoryResult{}, res.Error
	}
	rowsAffected := r.db.First(&employee, id).RowsAffected
	return usecases.OneRepositoryResult{
		Employee: employee,
		MetaData: usecases.RepositoryResultMetaData{
			RowsAffected: rowsAffected,
		},
	}, nil
}

func (r Repository) FindAll(ctx context.Context) (usecases.MultipleRepositoryResult, error) {
	var employees []entities.Employee
	res := r.db.Find(&employees)
	if res.Error != nil {
		return usecases.MultipleRepositoryResult{}, nil
	}
	return usecases.MultipleRepositoryResult{
		Employees: employees,
	}, nil
}

func (r Repository) Create(ctx context.Context, employee entities.Employee) (usecases.OneRepositoryResult, error) {
	tx, err := internal.GetDBTransaction(ctx)
	if err != nil {
		return usecases.OneRepositoryResult{}, err
	}

	if err := tx.Create(&employee).Error; err != nil {
		return usecases.OneRepositoryResult{}, err
	}
	rowsAffected := tx.RowsAffected

	var createdEmployee entities.Employee
	// TODO What happens if multiple clients create employees at the same time?
	if err := tx.Last(&createdEmployee).Error; err != nil {
		return usecases.OneRepositoryResult{}, err
	}
	return usecases.OneRepositoryResult{
		Employee: employee,
		MetaData: usecases.RepositoryResultMetaData{
			RowsAffected: rowsAffected,
		},
	}, nil
}

func (r Repository) UpdateIfExists(ctx context.Context, employee entities.Employee) (usecases.RepositoryResultMetaData, error) {
	tx, err := internal.GetDBTransaction(ctx)
	if err != nil {
		return usecases.RepositoryResultMetaData{}, err
	}
	res := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&entities.Employee{}, employee.Id)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return usecases.RepositoryResultMetaData{}, res.Error
	}
	rowsAffected := res.RowsAffected
	if rowsAffected == 0 {
		return usecases.RepositoryResultMetaData{
			RowsAffected: rowsAffected,
		}, nil
	}
	rowsAffected = tx.Save(&employee).RowsAffected
	return usecases.RepositoryResultMetaData{
		RowsAffected: rowsAffected,
	}, nil
}

func (r Repository) Delete(ctx context.Context, id int64) (usecases.RepositoryResultMetaData, error) {
	res := r.db.Delete(&entities.Employee{}, id)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return usecases.RepositoryResultMetaData{}, res.Error
	}
	rowsAffected := res.RowsAffected
	return usecases.RepositoryResultMetaData{
		RowsAffected: rowsAffected,
	}, nil
}
