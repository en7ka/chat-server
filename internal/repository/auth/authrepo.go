package auth

import (
	"github.com/en7ka/auth/pkg/auth_v1"
)

type RepoAccess struct {
	client authv1.AuthApiClient
}

func NewRepoAccess(client authv1.AuthApiClient) *RepoAccess {
	return &RepoAccess{client: client}
}
