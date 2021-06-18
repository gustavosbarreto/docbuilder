package user

import (
	"context"
	"errors"

	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/shellhub-io/shellhub/pkg/validator"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrConflict     = errors.New("conflict")
	ErrBadRequest   = errors.New("bad request")
)

type Service interface {
	UpdateDataUser(ctx context.Context, data *models.User, id string) ([]validator.InvalidField, error)
	UpdatePasswordUser(ctx context.Context, currentPassword, newPassword, id string) error
}

type service struct {
	store store.Store
}

func NewService(store store.Store) Service {
	return &service{store}
}

func (s *service) UpdateDataUser(ctx context.Context, data *models.User, id string) ([]validator.InvalidField, error) {
	var invalid []validator.InvalidField

	if _, _, err := s.store.UserGetByID(ctx, id, false); err != nil {
		return invalid, err
	}

	if invalidFields, err := validator.CheckValidation(data); err != nil {
		return invalidFields, ErrBadRequest
	}

	var checkUsername, checkEmail bool

	if user, err := s.store.UserGetByUsername(ctx, data.Username); err == nil && user.ID != id {
		checkUsername = true
		invalid = append(invalid, validator.InvalidField{"username", "conflict", "", ""})
	}

	if user, err := s.store.UserGetByEmail(ctx, data.Email); err == nil && user.ID != id {
		checkEmail = true
		invalid = append(invalid, validator.InvalidField{"email", "conflict", "", ""})
	}

	if checkUsername || checkEmail {
		return invalid, ErrConflict
	}

	return invalid, s.store.UserUpdateData(ctx, data, id)
}

func (s *service) UpdatePasswordUser(ctx context.Context, currentPassword, newPassword, id string) error {
	user, _, err := s.store.UserGetByID(ctx, id, false)
	if err != nil {
		return err
	}

	if user.Password == currentPassword {
		return s.store.UserUpdatePassword(ctx, newPassword, id)
	}

	return ErrUnauthorized
}
