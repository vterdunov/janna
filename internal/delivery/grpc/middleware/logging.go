package middleware

import (
	"path"
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

		service := path.Dir(info.FullMethod)[1:]
		method := path.Base(info.FullMethod)

		l.Info().
			Str("request_id", reqID).
			Str("service", service).
			Str("method", method).
			Msg("Calling endpoint")

		resp, err := handler(ctx, req)

		defer func(begin time.Time) {
			l.Info().
				Str("request_id", reqID).
				Str("service", service).
				Str("method", method).
				Str("took", time.Since(begin).String()).
				Err(err).
				Msg("Called endpoint")
		}(time.Now())

		return resp, err
	}
}
