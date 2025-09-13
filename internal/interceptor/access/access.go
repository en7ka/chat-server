package access

import (
	"context"

	"google.golang.org/grpc"
)

func (a AuthInterceptor) AccessInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := a.AccessService.Access(ctx, info.FullMethod); err != nil {
		return nil, err
	}

	return handler(ctx, req)
}
