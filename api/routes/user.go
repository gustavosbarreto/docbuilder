package routes

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/shellhub-io/shellhub/api/apicontext"
	"github.com/shellhub-io/shellhub/api/user"
	"github.com/shellhub-io/shellhub/pkg/models"
)

const (
	UpdateUserDataURL     = "/users/:id/data"
	UpdateUserPasswordURL = "/users/:id/password" //nolint:gosec
)

func UpdateUserData(c apicontext.Context) error {
	var req models.User

	if err := c.Bind(&req); err != nil {
		return err
	}

	ID := c.Param("id")

	svc := user.NewService(c.Store())

	if invalidFields, err := svc.UpdateDataUser(c.Ctx(), &req, ID); err != nil {
		switch {
		case err == user.ErrBadRequest:
			return c.JSON(http.StatusBadRequest, invalidFields)
		case err == user.ErrConflict:
			return c.JSON(http.StatusConflict, invalidFields)
		default:
			return err
		}
	}

	return c.JSON(http.StatusOK, nil)
}

func UpdateUserPassword(c apicontext.Context) error {
	var req struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	if err := c.Bind(&req); err != nil {
		return err
	}

	ID := c.Param("id")

	if req.CurrentPassword != "" {
		sum := sha256.Sum256([]byte(req.CurrentPassword))
		sumByte := sum[:]
		req.CurrentPassword = hex.EncodeToString(sumByte)
	}

	if req.NewPassword != "" {
		sum := sha256.Sum256([]byte(req.NewPassword))
		sumByte := sum[:]
		req.NewPassword = hex.EncodeToString(sumByte)
	}

	svc := user.NewService(c.Store())

	if err := svc.UpdatePasswordUser(c.Ctx(), req.CurrentPassword, req.NewPassword, ID); err != nil {
		switch {
		case err == user.ErrUnauthorized:
			return c.JSON(http.StatusForbidden, nil)
		default:
			return err
		}
	}

	return c.JSON(http.StatusOK, nil)
}
