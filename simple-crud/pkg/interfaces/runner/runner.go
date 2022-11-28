package runner

import (
	employee_controller "example/simple-crud/pkg/interfaces/controllers/web/employee"
	db_gateway "example/simple-crud/pkg/interfaces/gateways/db"
	employee_db "example/simple-crud/pkg/interfaces/gateways/db/employee"
	employee_presenter "example/simple-crud/pkg/interfaces/presenters/web/employee"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run() {

	dbHandle := db_gateway.OpenDb()
	tx := db_gateway.NewTransaction(dbHandle)

	// employees
	employeeController := employee_controller.NewController(
		employee_db.NewRepository(dbHandle), tx, employee_presenter.Factory{},
	)

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.GET("/employees/:id", employeeController.Get)
	e.GET("/employees", employeeController.GetAll)
	e.POST("/employees", employeeController.Post)
	e.PUT("/employees/:id", employeeController.Put)
	e.DELETE("/employees/:id", employeeController.Delete)

	serverPort := os.Getenv("SERVER_PORT")
	e.Logger.Fatal(e.Start(":" + serverPort))
}
