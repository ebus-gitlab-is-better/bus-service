package data

import (
	"bus-service/internal/conf"
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/go-kratos/kratos/v2/log"
)

type KeycloakAPI struct {
	client       *gocloak.GoCloak
	logger       *log.Helper
	clientId     string
	clientSecret string
	realm        string
	username     string
	password     string
}

func NewKeyCloakAPI(conf *conf.Data, client *gocloak.GoCloak, logger log.Logger) *KeycloakAPI {
	return &KeycloakAPI{
		client:       client,
		logger:       log.NewHelper(logger),
		clientId:     conf.Keycloak.ClientId,
		clientSecret: conf.Keycloak.ClientSecret,
		realm:        conf.Keycloak.Realm,
		username:     conf.Keycloak.Username,
		password:     conf.Keycloak.Password,
	}
}

func (api *KeycloakAPI) CheckToken(accessToken string) (*gocloak.IntroSpectTokenResult, error) {
	return api.client.RetrospectToken(
		context.TODO(),
		accessToken,
		api.clientId,
		api.clientSecret,
		api.realm)
}

func (api *KeycloakAPI) GetUserInfo(accessToken string) (*gocloak.UserInfo, error) {
	return api.client.GetUserInfo(
		context.TODO(),
		accessToken,
		api.realm)
}

func (api *KeycloakAPI) GetUserByID(userId string) (*gocloak.User, error) {
	token, err := api.client.LoginClient(context.TODO(),
		api.clientId,
		api.clientSecret,
		api.realm)
	if err != nil {
		return nil, err
	}
	return api.client.GetUserByID(
		context.TODO(),
		token.AccessToken,
		api.realm,
		userId,
	)
}

func (api *KeycloakAPI) GetDrivers(roleName string) ([]*gocloak.User, error) {
	token, err := api.client.LoginAdmin(
		context.TODO(),
		api.username,
		api.password,
		api.realm)
	if err != nil {
		return nil, err
	}

	return api.client.GetUsersByRoleName(
		context.TODO(),
		token.AccessToken,
		api.realm,
		roleName,
		gocloak.GetUsersByRoleParams{},
	)
}
