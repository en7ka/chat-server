package access

import (
	servinf "github.com/en7ka/chat-server/internal/service/servinterface"
)

type AuthInterceptor struct {
	AccessService servinf.Access
}

func NewAuthInterceptor(accessService servinf.Access) *AuthInterceptor {
	return &AuthInterceptor{
		AccessService: accessService,
	}
}
