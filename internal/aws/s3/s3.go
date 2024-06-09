//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type GetObjectAPI interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}
