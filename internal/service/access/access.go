package access

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	authPref = "Bearer "
)

// Access метод для проверки доступа пользователя
func (s *ServiceAcc) Access(ctx context.Context, path string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("missing metadata")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return fmt.Errorf("missing authorization header")
	}

	if !strings.HasPrefix(authHeader[0], authPref) {
		return fmt.Errorf("invalid authorization header")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPref)
	mod := metadata.New(map[string]string{"Authorization": "Bearer " + accessToken})
	clientCtx := metadata.NewOutgoingContext(ctx, mod)
	fmt.Println(path)

	return s.accessRepository.Access(clientCtx, path)
}
