//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

package source

import (
	"context"
)

type Source[T any] interface {
	Generate(ctx context.Context) (T, error)
}
