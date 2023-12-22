package observability

import (
	"context"
	"time"
)

// RepeatEvery will call the f every t
func RepeatEvery(ctx context.Context, f func(), t time.Duration) {
	timer := time.NewTicker(t)
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			go f()
		}
	}
}
