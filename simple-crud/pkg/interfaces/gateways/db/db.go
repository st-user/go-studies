package db

import (
	"context"
	"example/simple-crud/pkg/interfaces/gateways/db/internal"
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenDb() *gorm.DB {

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

type transaction struct {
	db *gorm.DB
}

func NewTransaction(db *gorm.DB) transaction {
	return transaction{
		db: db,
	}
}

func (t transaction) DoInTx(ctx context.Context, f func(ctx context.Context) (interface{}, error)) (interface{}, error) {

	var ret interface{}
	err := t.db.Transaction(func(tx *gorm.DB) error {
		ctx := context.WithValue(ctx, internal.TxKey{}, tx)
		txResult, err := f(ctx)
		if err != nil {
			return err
		}
		ret = txResult
		return nil
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
}
