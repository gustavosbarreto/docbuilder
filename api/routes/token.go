package routes

import (
	"net/http"

	"github.com/shellhub-io/shellhub/api/apicontext"
	"github.com/shellhub-io/shellhub/api/apierr"
	"github.com/shellhub-io/shellhub/api/authsvc"
	"github.com/shellhub-io/shellhub/api/token"
	"github.com/shellhub-io/shellhub/pkg/models"
)

const (
	ListTokenURL   = "/tokens"
	CreateTokenURL = "/tokens"
	GetTokenURL    = "/tokens/:id"
	DeleteTokenURL = "/tokens/:id/del"    //#nosec
	UpdateTokenURL = "/tokens/update/:id" //#nosec
)

func ListToken(c apicontext.Context) error {
	tokens, err := token.NewService(c.Store()).ListToken(c.Ctx(), c.Tenant().ID)
	if err != nil {
		return apierr.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, tokens)
}

func CreateToken(c apicontext.Context) error {
	if _, err := token.NewService(c.Store()).CreateToken(c.Ctx(), c.Tenant().ID); err != nil {
		return apierr.HandleError(c, err)
	}

	svc := authsvc.NewService(c.Store(), nil, nil)

	token, err := svc.AuthAPIToken(c.Ctx(), &models.APITokenAuthRequest{
		TenantID: c.Tenant().ID,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, token)
}

func GetToken(c apicontext.Context) error {
	token, err := token.NewService(c.Store()).GetToken(c.Ctx(), c.Tenant().ID, c.Param("id"))
	if err != nil {
		return apierr.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, token)
}

func DeleteToken(c apicontext.Context) error {
	if err := token.NewService(c.Store()).DeleteToken(c.Ctx(), c.Tenant().ID, c.Param("id")); err != nil {
		return apierr.HandleError(c, err)
	}

	return nil
}

func UpdateToken(c apicontext.Context) error {
	if err := token.NewService(c.Store()).UpdateToken(c.Ctx(), c.Tenant().ID, c.Param("id"), &models.APITokenUpdate{}); err != nil {
		return apierr.HandleError(c, err)
	}

	return nil
}
