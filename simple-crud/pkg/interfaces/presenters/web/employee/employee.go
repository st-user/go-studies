package employee

import (
	"context"
	entities "example/simple-crud/pkg/entities/employee"
	"example/simple-crud/pkg/interfaces/data"
	"example/simple-crud/pkg/interfaces/presenters/web/internal"
	usecases "example/simple-crud/pkg/usecases/employee"
	ucerrors "example/simple-crud/pkg/usecases/errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type EmployeeResponse struct {
	Id        int64           `json:"id"`
	Name      string          `json:"name"`
	StartDate data.SimpleDate `json:"startDate"`
}

func employeeResponse(e *entities.Employee) EmployeeResponse {
	return EmployeeResponse{
		Id:        e.Id,
		Name:      e.Name,
		StartDate: data.SimpleDate{Time: e.StartDate},
	}
}

type presenter struct {
	echoCtx echo.Context
}

type Factory struct {
}

func (f Factory) Create(ctx echo.Context) usecases.OutputPort {
	return presenter{echoCtx: ctx}

}

func (p presenter) WriteError(ctx context.Context, err error) error {
	switch uerr := err.(type) {
	case data.InputDataError:
		if uerr.HasDetail() {
			log.Warn(uerr.Error())
		}
		return p.echoCtx.JSON(http.StatusBadRequest,
			&internal.ErrorMessageResponse{Message: uerr.UserMessage()})
	case ucerrors.DataNotFoundError:
		if uerr.Error() != "" {
			log.Warn(uerr.Error())
		}
		return p.echoCtx.JSON(http.StatusNotFound,
			&internal.ErrorMessageResponse{Message: uerr.UserMessage()})
	default:
		log.Error(err)
		return p.echoCtx.JSON(http.StatusBadRequest,
			&internal.ErrorMessageResponse{Message: "Internal Server Error"})
	}
}

func (p presenter) Write(ctx context.Context, employee entities.Employee) error {
	return p.echoCtx.JSON(http.StatusOK, employeeResponse(&employee))
}

func (p presenter) WriteMultiple(ctx context.Context, employees []entities.Employee) error {
	responseJson := []EmployeeResponse{}
	for i := 0; i < len(employees); i++ {
		responseJson = append(responseJson, employeeResponse(&employees[i]))
	}
	return p.echoCtx.JSON(http.StatusOK, responseJson)
}

func (p presenter) Created(ctx context.Context, employee entities.Employee) error {
	return p.echoCtx.JSON(http.StatusCreated, employeeResponse(&employee))
}

func (p presenter) Updated(ctx context.Context, employee entities.Employee) error {
	return p.echoCtx.JSON(http.StatusOK, employeeResponse(&employee))
}

func (p presenter) Deleted(ctx context.Context) error {
	return p.echoCtx.NoContent(http.StatusNoContent)
}
