//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

package destination

import "tail-time/internal/tale"

type Destination interface {
	Save(tale tale.Tale) error
}
