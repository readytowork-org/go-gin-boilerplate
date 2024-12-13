package utils

import (
	"context"
	"time"
)

func GetContext() context.Context {
	// Create a context.Background()
	ctx := context.Background()

	// Create a context.WithCancel() to create a cancellable context
	defer context.WithCancel(ctx)

	// Create a context.WithTimeout() to create a context with a timeout
	timeout := 5 * time.Second
	defer context.WithTimeout(ctx, timeout)

	// Create a context.WithDeadline() to create a context with a deadline
	deadline := time.Now().Add(10 * time.Second)
	defer context.WithDeadline(ctx, deadline)

	return ctx
}
