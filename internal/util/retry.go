package util

import (
	"context"
	"fmt"
)

func Retry(ctx context.Context, name string, n int, fn func() error) error {
	for i := 1; i <= n; i++ {
		if fn() != nil {
			GLog.Warning(ctx, "retry %s failed at %d", name, i)
		} else {
			return nil
		}
	}
	return fmt.Errorf("%s retried %d times, but failed finally", name, n)
}
