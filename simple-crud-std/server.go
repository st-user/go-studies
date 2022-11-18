package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"simple-crud-std/controller"
	"simple-crud-std/data/dao"

	"github.com/go-sql-driver/mysql"
)

func createDBHandle() *sql.DB {
	cfg := mysql.Config{
		User:      os.Getenv("DB_USER"),
		Passwd:    os.Getenv("DB_PASS"),
		Net:       "tcp",
		Addr:      os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		DBName:    os.Getenv("DB_NAME"),
		ParseTime: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {

	db := createDBHandle()

	accessor := dao.NewEmployeeDAO(db)
	controller := controller.NewEmployeeController(accessor)
	http.HandleFunc("/employees/", controller.Handle())

	http.ListenAndServe(":8080", nil)
}
