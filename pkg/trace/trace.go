package trace

import (
	"context"

	"github.com/flamego/flamego"
	"github.com/google/uuid"
)

type contextKey struct{}

var activeSpanKey = contextKey{}

func Context(ctx context.Context) context.Context {
	x := uuid.New().String()
	return context.WithValue(ctx, activeSpanKey, x)
}

func Trace(ctx context.Context) string {
	traceValue := ctx.Value(activeSpanKey)
	if trace, ok := traceValue.(string); ok {
		return trace
	}
	return ""
}

func Tracer() flamego.ContextInvoker {
	return func(ctx flamego.Context) {
		r := ctx.Request().Clone(Context(context.Background()))
		ctx.Map(r)
	}
}
