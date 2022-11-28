package employee

import (
	"context"
	"fmt"

	entities "example/simple-crud/pkg/entities/employee"
	"example/simple-crud/pkg/usecases/errors"
	"example/simple-crud/pkg/usecases/transactions"
)

type InputPort interface {
	ReadByID(ctx context.Context, id int64) error
	ReadAll(ctx context.Context) error
	Create(ctx context.Context, employee entities.Employee) error
	Update(ctx context.Context, employee entities.Employee) error
	Delete(ctx context.Context, id int64) error
	HandleError(ctx context.Context, err error) error
}

type OutputPort interface {
	WriteError(ctx context.Context, err error) error
	Write(ctx context.Context, employee entities.Employee) error
	WriteMultiple(ctx context.Context, employees []entities.Employee) error
	Created(ctx context.Context, employee entities.Employee) error
	Updated(ctx context.Context, employee entities.Employee) error
	Deleted(ctx context.Context) error
}

type Repository interface {
	FindByID(ctx context.Context, id int64) (OneRepositoryResult, error)
	FindAll(ctx context.Context) (MultipleRepositoryResult, error)
	Create(ctx context.Context, employee entities.Employee) (OneRepositoryResult, error)
	UpdateIfExists(ctx context.Context, employee entities.Employee) (RepositoryResultMetaData, error)
	Delete(ctx context.Context, id int64) (RepositoryResultMetaData, error)
}

type RepositoryResultMetaData struct {
	RowsAffected int64
}

type OneRepositoryResult struct {
	Employee entities.Employee
	MetaData RepositoryResultMetaData
}

type MultipleRepositoryResult struct {
	Employees []entities.Employee
}

type interactor struct {
	repository Repository
	tx         transactions.Transaction
	outputPort OutputPort
}

func NewInputPort(repository Repository, tx transactions.Transaction, outputPort OutputPort) interactor {
	return interactor{
		repository: repository,
		tx:         tx,
		outputPort: outputPort,
	}
}

func (it interactor) ReadByID(ctx context.Context, id int64) error {
	ret, err := it.repository.FindByID(ctx, id)
	if err != nil {
		return it.outputPort.WriteError(ctx, errors.NewUnexpectedError(
			fmt.Sprintf("employee#GetById(%d)", id),
			err))
	}
	if ret.MetaData.RowsAffected == 0 {
		return it.outputPort.WriteError(ctx, errors.NewDataNotFoundErrorWithLabel("Employee"))
	}
	return it.outputPort.Write(ctx, ret.Employee)
}

func (it interactor) ReadAll(ctx context.Context) error {
	ret, err := it.repository.FindAll(ctx)
	if err != nil {
		return it.outputPort.WriteError(ctx, errors.NewUnexpectedError("employee#GetAll", err))
	}
	return it.outputPort.WriteMultiple(ctx, ret.Employees)
}

func (it interactor) Create(ctx context.Context, employee entities.Employee) error {
	ret, err := it.tx.DoInTx(ctx, func(ctx context.Context) (interface{}, error) {
		ret, err := it.repository.Create(ctx, employee)
		if err != nil {
			return nil, errors.NewUnexpectedError("employee#Create.inTx", err)
		}

		return ret, nil
	})
	if err != nil {
		return it.outputPort.WriteError(ctx, errors.NewUnexpectedError("employee#Create.Tx", err))
	}
	if e, ok := ret.(OneRepositoryResult); !ok {
		return it.outputPort.WriteError(ctx, errors.NewUnexpectedError("employee#Create.type", err))
	} else {
		return it.outputPort.Created(ctx, e.Employee)
	}
}

func (it interactor) Update(ctx context.Context, employee entities.Employee) error {
	ret, err := it.tx.DoInTx(ctx, func(ctx context.Context) (interface{}, error) {
		ret, err := it.repository.UpdateIfExists(ctx, employee)
		if err != nil {
			return nil, errors.NewUnexpectedError("employee#Update.inTx", err)
		}
		return ret, nil
	})
	if err != nil {
		return it.outputPort.WriteError(ctx, errors.NewUnexpectedError("employee#Update.Tx", err))
	}
	if ret, ok := ret.(RepositoryResultMetaData); !ok {
		return it.outputPort.WriteError(ctx, errors.NewUnexpectedError("employee#Update.type", err))
	} else {
		if ret.RowsAffected == 0 {
			return it.outputPort.WriteError(ctx, errors.NewDataNotFoundErrorWithLabel("Employee"))
		}
		return it.outputPort.Updated(ctx, employee)
	}

}

func (it interactor) Delete(ctx context.Context, id int64) error {
	_, err := it.tx.DoInTx(ctx, func(ctx context.Context) (interface{}, error) {
		_, err := it.repository.Delete(ctx, id)
		if err != nil {
			return nil, errors.NewUnexpectedError("employee#Delete.inTx", err)
		}
		return nil, nil
	})
	if err != nil {
		return it.outputPort.WriteError(ctx, errors.NewUnexpectedError("employee#Delete.Tx", err))
	}
	return it.outputPort.Deleted(ctx)
}

func (it interactor) HandleError(ctx context.Context, err error) error {
	return it.outputPort.WriteError(ctx, err)
}
