//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

package destination

type Destination interface {
	Save(tale string) error
}
