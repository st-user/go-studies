package internal

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type TxKey struct{}

func GetDBTransaction(ctx context.Context) (*gorm.DB, error) {
	if tx, ok := ctx.Value(TxKey{}).(*gorm.DB); !ok {
		return nil, errors.New("must be done in a transaction")
	} else {
		return tx, nil
	}
}
