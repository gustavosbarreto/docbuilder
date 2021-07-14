package store

import (
	"context"

	"github.com/shellhub-io/shellhub/pkg/api/paginator"
	"github.com/shellhub-io/shellhub/pkg/models"
)

type NamespaceStore interface {
	NamespaceList(ctx context.Context, pagination paginator.Query, filters []models.Filter, export bool) ([]models.Namespace, int, error)
	NamespaceGet(ctx context.Context, tenantID string) (*models.Namespace, error)
	NamespaceGetByName(ctx context.Context, tenantID string) (*models.Namespace, error)
	NamespaceCreate(ctx context.Context, tenantID *models.Namespace) (*models.Namespace, error)
	NamespaceRename(ctx context.Context, tenantID, name string) (*models.Namespace, error)
	NamespaceUpdate(ctx context.Context, tenantID string, namespace *models.Namespace) error
	NamespaceDelete(ctx context.Context, tenantID string) error
	NamespaceAddMember(ctx context.Context, tenantID, ID string) (*models.Namespace, error)
	NamespaceRemoveMember(ctx context.Context, tenantID, ID string) (*models.Namespace, error)
	NamespaceGetFirst(ctx context.Context, ID string) (*models.Namespace, error)
	NamespaceSetSessionRecord(ctx context.Context, sessionRecord bool, tenantID string) error
	NamespaceGetSessionRecord(ctx context.Context, tenantID string) (bool, error)
	NamespaceSetWebhook(ctx context.Context, tenant, url string) (*models.Namespace, error)
	NamespaceSetWebhookStatus(ctx context.Context, tenant string, active bool) (*models.Namespace, error)
}
