package connections

import (
	"context"

	"fampay-assignment/config"
	"fampay-assignment/utils"
)

var (
	ctx context.Context = context.Background()
)

func GetContext() context.Context {
	return utils.GetContextWithTimeout(ctx, config.QUERY_TIMEOUT)
}
