package middleware

import (
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func LoggingInterceptor(l zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		reqID := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			reqID = md["request_id"][0]
		}

		l.Info().
			Str("request_id", reqID).
			Str("method", info.FullMethod).
			Msg("Calling endpoint")

		defer func(begin time.Time) {
			l.Info().
				Str("request_id", reqID).
				Str("method", info.FullMethod).
				Str("took", time.Since(begin).String()).
				Msg("Called endpoint")
		}(time.Now())

		// TODO: log errors

		return handler(ctx, req)
	}
}
