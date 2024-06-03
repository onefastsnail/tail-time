//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

package destination

type Destination[T any] interface {
	Save(tale T) error
}
