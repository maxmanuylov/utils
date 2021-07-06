package utils

import (
	"context"
	"io"
)

func Close(closer io.Closer) {
	_ = closer.Close()
}

func Defer(ctx context.Context, action func()) {
	if done := ctx.Done(); done != nil {
		go func() {
			<-done
			action()
		}()
	}
}

func DeferClose(ctx context.Context, closer io.Closer) {
	Defer(ctx, func() {
		Close(closer)
	})
}
