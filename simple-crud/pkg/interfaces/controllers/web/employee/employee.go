package employee

import (
	entities "example/simple-crud/pkg/entities/employee"
	"example/simple-crud/pkg/interfaces/data"
	usecases "example/simple-crud/pkg/usecases/employee"
	userTx "example/simple-crud/pkg/usecases/transactions"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

type EmployeeRequest struct {
	Name      string          `json:"name"`
	StartDate data.SimpleDate `json:"startDate"`
}

func (e *EmployeeRequest) employee() entities.Employee {
	return entities.Employee{
		Name:      e.Name,
		StartDate: e.StartDate.Time,
	}
}

type controller struct {
	repository        usecases.Repository
	tx                userTx.Transaction
	outputPortFactory OutputPortFactory
}

type OutputPortFactory interface {
	Create(ctx echo.Context) usecases.OutputPort
}

func NewController(repository usecases.Repository, tx userTx.Transaction, outputPortFactory OutputPortFactory) controller {
	return controller{
		repository:        repository,
		tx:                tx,
		outputPortFactory: outputPortFactory,
	}
}

func (c controller) newInputPort(ctx echo.Context) usecases.InputPort {
	return usecases.NewInputPort(
		c.repository, c.tx, c.outputPortFactory.Create(ctx),
	)
}

func (c controller) Get(ctx echo.Context) error {
	inputPort := c.newInputPort(ctx)

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return inputPort.HandleError(ctx.Request().Context(), data.NewInputDataError(
			fmt.Sprintf("Invalid id value: %s", idStr), nil,
		))
	}

	return inputPort.ReadByID(ctx.Request().Context(), int64(id))
}

func (c controller) GetAll(ctx echo.Context) error {
	inputPort := c.newInputPort(ctx)
	return inputPort.ReadAll(ctx.Request().Context())
}

func (c controller) Post(ctx echo.Context) error {
	inputPort := c.newInputPort(ctx)

	body := new(EmployeeRequest)
	if err := ctx.Bind(body); err != nil {
		return inputPort.HandleError(ctx.Request().Context(), data.NewInputDataError(
			"Can not parse the request body", err,
		))
	}

	return inputPort.Create(ctx.Request().Context(), body.employee())
}

func (c controller) Put(ctx echo.Context) error {
	inputPort := c.newInputPort(ctx)

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return inputPort.HandleError(ctx.Request().Context(), data.NewInputDataError(
			fmt.Sprintf("Invalid id value: %s", idStr), nil,
		))
	}

	body := new(EmployeeRequest)
	if err := ctx.Bind(body); err != nil {
		return inputPort.HandleError(ctx.Request().Context(), data.NewInputDataError(
			"Can not parse the request body", err,
		))
	}
	newEmployee := body.employee()
	newEmployee.Id = int64(id)

	return inputPort.Update(ctx.Request().Context(), newEmployee)
}

func (c controller) Delete(ctx echo.Context) error {
	inputPort := c.newInputPort(ctx)

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return inputPort.HandleError(ctx.Request().Context(), data.NewInputDataError(
			fmt.Sprintf("Invalid id value: %s", idStr), nil,
		))
	}

	return inputPort.Delete(ctx.Request().Context(), int64(id))
}
