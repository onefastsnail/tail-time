//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

package source

import (
	"context"

	"tail-time/internal/tale"
)

type Source interface {
	Generate(ctx context.Context) (tale.Tale, error)
}
