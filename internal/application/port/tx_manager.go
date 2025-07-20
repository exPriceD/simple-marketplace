package port

import "context"

// TxManager — задел на будущее (транзакции при нескольких репозиториях).
type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}
