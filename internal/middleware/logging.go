package middleware

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
)

func LoggingInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		in interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, in)
		if err != nil {
			log.Error("Request failed", slog.String("method", info.FullMethod), slog.Any("error", err))
		} else {
			log.Info("Request succeeded", slog.String("method", info.FullMethod))
		}
		return resp, err
	}
}
