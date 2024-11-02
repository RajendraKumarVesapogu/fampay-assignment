package utils

import (
	"context"
	"time"	
)


func GetPaginationOffset(page int, size int) int {
	return (page - 1) * size
}

func GetContextWithTimeout(ctx context.Context, timeout time.Duration) context.Context {
	Ctx, cancel := context.WithCancel(ctx)

	timer := time.NewTimer(timeout)

	go func() {
		defer cancel()
		<-timer.C
	}()

	return Ctx
}
