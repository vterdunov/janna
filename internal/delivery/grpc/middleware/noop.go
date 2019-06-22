package middleware

import (
	"context"

	"google.golang.org/grpc"
)

func NoopInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// md, ok := metadata.FromIncomingContext(ctx)
	// if !ok {
	// 	fmt.Println("Error")
	// }
	// spew.Dump(md)
	return handler(ctx, req)
}
