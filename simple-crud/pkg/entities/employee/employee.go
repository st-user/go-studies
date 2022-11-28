package employee

import "time"

type Employee struct {
	Id        int64
	Name      string
	StartDate time.Time
}
