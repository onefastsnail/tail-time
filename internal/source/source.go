//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

package source

import "context"

type Source interface {
	Generate(ctx context.Context) (string, error)
}
