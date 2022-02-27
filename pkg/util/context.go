package util

import (
	"context"
	"time"
)

type HandlerWithContext func(ctx context.Context) (interface{}, error)

func WithTimeout(h HandlerWithContext, timeout int64) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	r, err := h(ctx)

	return r, err
}

