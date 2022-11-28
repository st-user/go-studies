package transactions

import "context"

type Transaction interface {
	DoInTx(ctx context.Context, f func(ctx context.Context) (interface{}, error)) (interface{}, error)
}
