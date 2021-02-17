package token

import (
	"context"
	"testing"

	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/api/store/mocks"
	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/stretchr/testify/assert"
)

type Data struct {
	Mock      *mocks.Store
	Service   Service
	Context   context.Context
	Namespace *models.Namespace
	Token     *models.Token
	err       error
}

func TestListToken(t *testing.T) {
	data := initData(t)

	data.Mock.On("TokenListAPIToken", data.Context, data.Namespace.TenantID).Return([]models.Token{}, nil).Once()

	_, err := data.Service.ListToken(data.Context, data.Namespace.TenantID)
	assert.NoError(t, err)

	data.Mock.AssertExpectations(t)
}

func TestCreateToken(t *testing.T) {
	data := initData(t)

	assert.NoError(t, data.err)
}

func TestGetToken(t *testing.T) {
	data := initData(t)

	data.Mock.On("TokenGetAPIToken", data.Context, data.Namespace.TenantID, data.Token.ID).Return(&models.Token{}, nil).Once()

	_, err := data.Service.GetToken(data.Context, data.Namespace.TenantID, data.Token.ID)
	assert.NoError(t, err)

	data.Mock.AssertExpectations(t)
}

func TestDeleteToken(t *testing.T) {
	data := initData(t)

	data.Mock.On("TokenDeleteAPIToken", data.Context, data.Namespace.TenantID, data.Token.ID).Return(nil).Once()

	err := data.Service.DeleteToken(data.Context, data.Namespace.TenantID, data.Token.ID)
	assert.NoError(t, err)

	data.Mock.AssertExpectations(t)
}

func TestUpdateToken(t *testing.T) {
	data := initData(t)

	req := &models.APITokenUpdate{
		TokenFields: models.TokenFields{ReadOnly: false},
	}

	data.Mock.On("TokenUpdateAPIToken", data.Context, data.Namespace.TenantID, data.Token.ID, req).Return(nil, nil).Once()

	err := data.Service.UpdateToken(data.Context, data.Namespace.TenantID, data.Token.ID, req)
	assert.NoError(t, err)

	data.Mock.On("TokenGetAPIToken", data.Context, data.Namespace.TenantID, data.Token.ID).Return(&models.Token{}, nil).Once()

	returnedToken, err := data.Service.GetToken(data.Context, data.Namespace.TenantID, data.Token.ID)
	assert.NoError(t, err)
	assert.Equal(t, returnedToken.ReadOnly, false)

	data.Mock.AssertExpectations(t)
}

func initData(t *testing.T) (response Data) {
	t.Helper()
	mock := &mocks.Store{}

	ctx := context.TODO()

	svc := NewService(store.Store(mock))

	token := models.Token{
		ID:       "1",
		TenantID: "a736a52b-5777-4f92-b0b8-e359bf484713",
		ReadOnly: true,
	}

	namespace := &models.Namespace{Name: "group1", Owner: "hash1", TenantID: "a736a52b-5777-4f92-b0b8-e359bf484713", APITokens: []models.Token{}}

	mock.On("TokenCreateAPIToken", ctx, namespace.TenantID).Return(&token, nil).Once()

	createdToken, err := svc.CreateToken(ctx, namespace.TenantID)
	assert.NoError(t, err)

	return Data{
		Mock:      mock,
		Service:   svc,
		Context:   ctx,
		Namespace: namespace,
		Token:     createdToken,
		err:       nil,
	}
}
