package util

import (
	"context"
	"fmt"
)

func Retry(ctx context.Context, name string, n int, fn func() error) error {
	for i := 1; i <= n; i++ {
		if err := fn(); err != nil {
			GLog.Warningf(ctx, "retry %s failed at %d, because %s", name, i, err.Error())
		} else {
			return nil
		}
	}
	return fmt.Errorf("%s retried %d times, but failed finally", name, n)
}
