package store

import (
	"context"

	"github.com/shellhub-io/shellhub/pkg/models"
)

type TokenStore interface {
	TokenListAPIToken(ctx context.Context, tenantID string) ([]models.Token, error)
	TokenCreateAPIToken(ctx context.Context, tenantID string) (*models.Token, error)
	TokenGetAPIToken(ctx context.Context, tenantID string, ID string) (*models.Token, error)
	TokenDeleteAPIToken(ctx context.Context, tenantID string, ID string) error
	TokenUpdateAPIToken(ctx context.Context, tenantID string, ID string, token *models.APITokenUpdate) error
}
