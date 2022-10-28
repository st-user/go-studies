package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Employee struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"startDate"`
}

type EmployeeRequest struct {
	Name      string     `json:"name"`
	StartDate SimpleDate `json:"startDate"`
}

func (e *EmployeeRequest) toEmployee() Employee {
	return Employee{
		Name:      e.Name,
		StartDate: e.StartDate.Time,
	}
}

type EmployeeResponse struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name"`
	StartDate SimpleDate `json:"startDate"`
}

func toEmployeeResponse(e *Employee) EmployeeResponse {
	return EmployeeResponse{
		Id:        e.Id,
		Name:      e.Name,
		StartDate: SimpleDate{e.StartDate},
	}
}

type ErrorMessageResponse struct {
	Message string `json:"message"`
}

type SimpleDate struct {
	time.Time
}

func (d SimpleDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Time.Format("2006-01-02") + `"`), nil
}

func (d *SimpleDate) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	parsedTime, err := time.ParseInLocation(`"2006-01-02"`, string(data), time.UTC)
	if err != nil {
		return err
	}

	*d = SimpleDate{parsedTime}
	return err
}

type NotFoundError struct {
	msg string
}

func (e *NotFoundError) Error() string {
	return e.msg
}

func openDb() *gorm.DB {

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// dsn := "root:my-secret-pw@tcp(127.0.0.1:3307)/testdb?charset=utf8mb4"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	dbDebugStr := os.Getenv("DB_DEBUG")
	dbDebug, err := strconv.ParseBool(dbDebugStr)

	if err != nil {
		fmt.Printf("%s\n", err)
		dbDebug = false
	}

	if dbDebug {
		return db.Debug()
	} else {
		return db
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	db := openDb()

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	/*
	 * GET /employees
	 */
	e.GET("/employees", func(c echo.Context) error {
		var employees []Employee
		db.Find(&employees)
		responseJson := []EmployeeResponse{}
		for i := 0; i < len(employees); i++ {
			responseJson = append(responseJson, toEmployeeResponse(&employees[i]))
		}
		return c.JSON(http.StatusOK, responseJson)
	})

	/*
	 * GET /employees/:id
	 */
	e.GET("/employees/:id", func(c echo.Context) error {

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Warn(err)
			return c.JSON(http.StatusBadRequest,
				&ErrorMessageResponse{fmt.Sprintf("Invalid id value: %s", idStr)})
		}

		var employee Employee
		rowsAffected := db.First(&employee, id).RowsAffected
		if rowsAffected == 0 {
			return c.JSON(http.StatusNotFound, &ErrorMessageResponse{"Employee Not Found"})
		}

		return c.JSON(http.StatusOK, toEmployeeResponse(&employee))
	})

	/*
	 * POST /employees
	 */
	e.POST("/employees", func(c echo.Context) error {

		body := new(EmployeeRequest)
		if err := c.Bind(body); err != nil {
			log.Warn(err)
			return c.JSON(http.StatusBadRequest,
				&ErrorMessageResponse{"Can not parse the request body"})
		}

		var responseEmployee Employee
		employee := body.toEmployee()
		if err := db.Transaction(func(tx *gorm.DB) error {

			if err := tx.Create(&employee).Error; err != nil {
				return err
			}
			// TODO What happens if multiple clients create employees at the same time?
			tx.Last(&responseEmployee)
			return nil

		}); err != nil {
			log.Error(err)
			return c.JSON(http.StatusInternalServerError,
				&ErrorMessageResponse{"Encountered an error while creating employee"})
		}

		return c.JSON(http.StatusCreated, toEmployeeResponse(&responseEmployee))
	})

	/*
	 * PUT /employees/:id
	 */
	e.PUT("/employees/:id", func(c echo.Context) error {

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Warn(err)
			return c.JSON(http.StatusBadRequest,
				&ErrorMessageResponse{fmt.Sprintf("Invalid id value: %s", idStr)})
		}

		body := new(EmployeeRequest)
		if err := c.Bind(body); err != nil {
			log.Warn(err)
			return c.JSON(http.StatusBadRequest,
				&ErrorMessageResponse{"Can not parse the request body"})
		}

		newEmployee := body.toEmployee()
		newEmployee.Id = int64(id)
		if err := db.Transaction(func(tx *gorm.DB) error {

			rowAffected := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&Employee{}, id).RowsAffected
			if rowAffected == 0 {
				return &NotFoundError{"Employee Not Found"}
			}

			tx.Save(&newEmployee)
			return nil

		}); err != nil {

			notFoundErr, ok := err.(*NotFoundError)
			if ok {
				return c.JSON(http.StatusNotFound, &ErrorMessageResponse{notFoundErr.msg})
			}
			log.Error(err)
			return c.JSON(http.StatusInternalServerError,
				&ErrorMessageResponse{"Encountered an error while updating employee"})

		}
		return c.JSON(http.StatusOK, toEmployeeResponse(&newEmployee))
	})

	/*
	 * DELETE /employees/:id
	 */
	e.DELETE("/employees/:id", func(c echo.Context) error {

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Warn(err)
			return c.JSON(http.StatusBadRequest,
				&ErrorMessageResponse{fmt.Sprintf("Invalid id value: %s", idStr)})
		}
		db.Delete(&Employee{}, id)

		return c.NoContent(http.StatusNoContent)
	})

	serverPort := os.Getenv("SERVER_PORT")
	e.Logger.Fatal(e.Start(":" + serverPort))
}
