package middleware

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

func LoggingInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		in interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Info("Request started", slog.String("method", info.FullMethod))

		t := time.Now()
		resp, err := handler(ctx, in)
		latency := time.Since(t)

		if err != nil {
			log.Error("Request failed",
				slog.String("method", info.FullMethod),
				slog.String("error", err.Error()),
				slog.Duration("latency", latency),
			)
		} else {
			log.Info("Request succeeded",
				slog.String("method", info.FullMethod),
				slog.Duration("latency", latency),
			)
		}
		return resp, err
	}
}
