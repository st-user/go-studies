package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"simple-crud-std/data/entity"
	"strconv"
	"time"
)

type EmployeeRequest struct {
	Name      string     `json:"name"`
	StartDate SimpleDate `json:"startDate"`
}

func (e *EmployeeRequest) toEmployee() entity.Employee {
	return entity.Employee{
		Name:      e.Name,
		StartDate: e.StartDate.Time,
	}
}

type EmployeeResponse struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	StartDate SimpleDate `json:"startDate"`
}

func toEmployeeResponse(e *entity.Employee) EmployeeResponse {
	return EmployeeResponse{
		ID:        e.ID,
		Name:      e.Name,
		StartDate: SimpleDate{e.StartDate},
	}
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

type requestValidationError struct {
	msg string
}

func (e *requestValidationError) Error() string {
	return e.msg
}

type EmployeeDataAccessor interface {
	FindAll() ([]entity.Employee, error)
	FindByID(id int64) (entity.Employee, error)
	Save(e entity.Employee) (entity.Employee, error)
	Update(e entity.Employee) error
	DeleteByID(id int64) error
}

type EmployeeController struct {
	accessor EmployeeDataAccessor
}

func NewEmployeeController(accessor EmployeeDataAccessor) EmployeeController {
	return EmployeeController{accessor}
}

// Handle creates a handler function (http.HandlerFunc)
// that handles '/employees' endpoints.
// This handler accepts GET, POST, PUT and DELETE methods.
func (ec EmployeeController) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/employees/"):]

		var resp interface{}
		var err error
		switch r.Method {
		case http.MethodGet, "":
			resp, err = ec.get(id)
		case http.MethodPost:
			resp, err = ec.post(r)
		case http.MethodPut:
			resp, err = ec.put(id, r)
		case http.MethodDelete:
			resp, err = ec.delete(id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if err != nil {
			switch err.(type) {
			case *entity.NotFoundError:
				fmt.Printf("%s\n", err)
				w.WriteHeader(http.StatusNotFound)
			case *requestValidationError:
				fmt.Printf("%s\n", err)
				w.WriteHeader(http.StatusBadRequest)
			default:
				fmt.Fprintf(os.Stderr, "%s\n", err)
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		if resp == nil {
			w.WriteHeader(http.StatusNoContent)
		} else {
			dec := json.NewEncoder(w)
			err = dec.Encode(resp)

			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

	}
}

func (ec EmployeeController) get(idStr string) (interface{}, error) {
	if idStr == "" {
		employees, err := ec.accessor.FindAll()
		if err != nil {
			return 0, err
		}
		resp := []EmployeeResponse{}
		for _, e := range employees {
			resp = append(resp, toEmployeeResponse(&e))
		}

		return resp, nil

	} else {

		id, err := stringIDtoInt(idStr)
		if err != nil {
			return 0, err
		}

		employee, err := ec.accessor.FindByID(id)
		if err != nil {
			return 0, err
		}

		employeeResp := toEmployeeResponse(&employee)
		return &employeeResp, nil

	}
}

func (ec EmployeeController) post(r *http.Request) (interface{}, error) {

	var body EmployeeRequest
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&body)
	if err != nil {
		return 0, err
	}

	newEmployee, err := ec.accessor.Save(body.toEmployee())
	if err != nil {
		return 0, err
	}

	employeeResp := toEmployeeResponse(&newEmployee)
	return &employeeResp, nil
}

func (ec EmployeeController) put(idStr string, r *http.Request) (interface{}, error) {

	var body EmployeeRequest
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&body)
	if err != nil {
		return 0, err
	}

	id, err := stringIDtoInt(idStr)
	if err != nil {
		return 0, err
	}

	employee := body.toEmployee()
	employee.ID = id

	err = ec.accessor.Update(employee)
	if err != nil {
		return 0, err
	}

	employeeResp := toEmployeeResponse(&employee)
	return &employeeResp, nil
}

func (ec EmployeeController) delete(idStr string) (interface{}, error) {

	id, err := stringIDtoInt(idStr)
	if err != nil {
		return 0, err
	}

	err = ec.accessor.DeleteByID(id)
	if err != nil {
		return 0, err
	}

	return nil, nil
}

func stringIDtoInt(idStr string) (int64, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, &requestValidationError{
			fmt.Sprintf("invalid ID value: %s", idStr),
		}
	}
	return int64(id), nil
}
