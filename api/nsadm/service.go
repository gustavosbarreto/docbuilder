package nsadm

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

	utils "github.com/shellhub-io/shellhub/api/pkg/namespace"
	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/pkg/api/paginator"
	"github.com/shellhub-io/shellhub/pkg/envs"
	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/shellhub-io/shellhub/pkg/uuid"
	"github.com/shellhub-io/shellhub/pkg/validator"
)

type Service interface {
	ListNamespaces(ctx context.Context, pagination paginator.Query, filterB64 string, export bool) ([]models.Namespace, int, error)
	CreateNamespace(ctx context.Context, namespace *models.Namespace, ownerUsername string) (*models.Namespace, error)
	GetNamespace(ctx context.Context, tenantID string) (*models.Namespace, error)
	DeleteNamespace(ctx context.Context, tenantID, ownerUsername string) error
	EditNamespace(ctx context.Context, tenantID, name, ownerUsername string) (*models.Namespace, error)
	AddNamespaceUser(ctx context.Context, tenantID, username, ownerUsername string) (*models.Namespace, error)
	RemoveNamespaceUser(ctx context.Context, tenantID, username, ownerUsername string) (*models.Namespace, error)
	ListMembers(ctx context.Context, tenantID string) ([]models.Member, error)
	EditSessionRecordStatus(ctx context.Context, status bool, tenant, ownerID string) error
	GetSessionRecord(ctx context.Context, tenant string) (bool, error)
	GetNamespaceByName(ctx context.Context, tenant string) (*models.Namespace, error)
	UpdateWebhook(ctx context.Context, url, tenant, owner string) ([]validator.InvalidField, *models.Namespace, error)
	SetWebhookStatus(ctx context.Context, active bool, tenant, owner string) (*models.Namespace, error)
}

type service struct {
	store store.Store
}

func NewService(store store.Store) Service {
	return &service{store}
}

func (s *service) ListNamespaces(ctx context.Context, pagination paginator.Query, filterB64 string, export bool) ([]models.Namespace, int, error) {
	raw, err := base64.StdEncoding.DecodeString(filterB64)
	if err != nil {
		return nil, 0, err
	}

	var filter []models.Filter

	if err := json.Unmarshal(raw, &filter); len(raw) > 0 && err != nil {
		return nil, 0, err
	}

	return s.store.NamespaceList(ctx, pagination, filter, export)
}

func (s *service) CreateNamespace(ctx context.Context, namespace *models.Namespace, ownerID string) (*models.Namespace, error) {
	user, _, err := s.store.UserGetByID(ctx, ownerID, false)
	if user == nil {
		return nil, ErrUnauthorized
	}

	if err != nil {
		return nil, err
	}

	ns := &models.Namespace{
		Name:     strings.ToLower(namespace.Name),
		Owner:    user.ID,
		Members:  []interface{}{user.ID},
		Settings: &models.NamespaceSettings{SessionRecord: true},
		TenantID: namespace.TenantID,
	}

	if _, err := validator.ValidateStruct(ns); err != nil {
		return nil, ErrInvalidFormat
	}

	if namespace.TenantID == "" {
		ns.TenantID = uuid.Generate()
	}

	// Set limits according to ShellHub instance type
	if envs.IsCloud() {
		// cloud free plan is limited only by the max of devices
		ns.MaxDevices = 3
	} else {
		// we don't set limits on enterprise and community instances
		ns.MaxDevices = -1
	}

	otherNamespace, err := s.store.NamespaceGetByName(ctx, ns.Name)
	if err != nil && err != store.ErrNoDocuments {
		return nil, err
	}

	if otherNamespace != nil {
		return nil, ErrConflictName
	}

	if _, err := s.store.NamespaceCreate(ctx, ns); err != nil {
		return nil, err
	}

	return ns, nil
}

func (s *service) GetNamespace(ctx context.Context, tenantID string) (*models.Namespace, error) {
	return s.store.NamespaceGet(ctx, tenantID)
}

func (s *service) DeleteNamespace(ctx context.Context, tenantID, ownerID string) error {
	if err := utils.IsNamespaceOwner(ctx, s.store, tenantID, ownerID); err != nil {
		return err
	}

	return s.store.NamespaceDelete(ctx, tenantID)
}

func (s *service) ListMembers(ctx context.Context, tenantID string) ([]models.Member, error) {
	ns, err := s.store.NamespaceGet(ctx, tenantID)
	if err == store.ErrNoDocuments {
		return nil, ErrNamespaceNotFound
	}

	if err != nil {
		return nil, err
	}

	members := []models.Member{}
	for _, memberID := range ns.Members {
		user, _, err := s.store.UserGetByID(ctx, memberID.(string), false)
		if err == store.ErrNoDocuments {
			return nil, ErrUserNotFound
		}

		if err != nil {
			return nil, err
		}

		member := models.Member{ID: memberID.(string), Name: user.Username}
		members = append(members, member)
	}

	return members, nil
}

func (s *service) EditNamespace(ctx context.Context, tenantID, name, owner string) (*models.Namespace, error) {
	if err := utils.IsNamespaceOwner(ctx, s.store, tenantID, owner); err != nil {
		return nil, err
	}

	ns, err := s.store.NamespaceGet(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	lowerName := strings.ToLower(name)
	if _, err := validator.ValidateStruct(&models.Namespace{
		Name: lowerName,
	}); err != nil {
		return nil, ErrInvalidFormat
	}

	if ns.Name == lowerName {
		return nil, ErrUnauthorized
	}

	return s.store.NamespaceRename(ctx, ns.TenantID, lowerName)
}

func (s *service) AddNamespaceUser(ctx context.Context, tenantID, username, ownerID string) (*models.Namespace, error) {
	if err := utils.IsNamespaceOwner(ctx, s.store, tenantID, ownerID); err != nil {
		return nil, err
	}

	user, err := s.store.UserGetByUsername(ctx, username)
	if err == store.ErrNoDocuments {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return s.store.NamespaceAddMember(ctx, tenantID, user.ID)
}

func (s *service) RemoveNamespaceUser(ctx context.Context, tenantID, username, ownerID string) (*models.Namespace, error) {
	if err := utils.IsNamespaceOwner(ctx, s.store, tenantID, ownerID); err != nil {
		return nil, err
	}

	user, err := s.store.UserGetByUsername(ctx, username)
	if err == store.ErrNoDocuments {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return s.store.NamespaceRemoveMember(ctx, tenantID, user.ID)
}

func (s *service) EditSessionRecordStatus(ctx context.Context, sessionRecord bool, tenant, ownerID string) error {
	if err := utils.IsNamespaceOwner(ctx, s.store, tenant, ownerID); err != nil {
		return err
	}

	return s.store.NamespaceSetSessionRecord(ctx, sessionRecord, tenant)
}

func (s *service) GetSessionRecord(ctx context.Context, tenant string) (bool, error) {
	if _, err := s.store.NamespaceGet(ctx, tenant); err != nil {
		if err == store.ErrNoDocuments {
			return false, ErrNamespaceNotFound
		}

		return false, err
	}

	return s.store.NamespaceGetSessionRecord(ctx, tenant)
}

func (s *service) GetNamespaceByName(ctx context.Context, namespace string) (*models.Namespace, error) {
	return s.store.NamespaceGetByName(ctx, namespace)
}

func (s *service) UpdateWebhook(ctx context.Context, url, tenant, owner string) ([]validator.InvalidField, *models.Namespace, error) {
	if err := utils.IsNamespaceOwner(ctx, s.store, tenant, owner); err != nil {
		return nil, nil, err
	}

	if invalidFields, err := validator.ValidateVar(url, "required,url"); err != nil {
		return invalidFields, nil, ErrBadRequest
	}
	ns, err := s.store.NamespaceSetWebhook(ctx, tenant, url)

	return nil, ns, err
}

func (s *service) SetWebhookStatus(ctx context.Context, active bool, tenant, owner string) (*models.Namespace, error) {
	if err := utils.IsNamespaceOwner(ctx, s.store, tenant, owner); err != nil {
		return nil, err
	}

	return s.store.NamespaceSetWebhookStatus(ctx, tenant, active)
}
