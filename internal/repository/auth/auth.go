package auth

import (
	"context"

	"github.com/en7ka/auth/pkg/auth_v1"
)

// Access вызываем сервис авторизации для проверки доступа.
func (r *RepoAccess) Access(ctx context.Context, path string) error {
	_, err := r.client.Check(ctx, &authv1.CheckRequest{EndpointAddress: path})
	return err
}
