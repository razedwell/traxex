package obs

import (
	"context"

	"go.opentelemetry.io/otel"
)

func InitTracer(ctx context.Context, serviceName string) (func(context.Context), error, error) {
	otel.SetTracerProvider()
}
