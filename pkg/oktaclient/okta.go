package oktaclient

import (
	"context"
	"fmt"

	"github.com/aserto-dev/aserto-idp-plugin-okta/pkg/config"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=okta.go -destination=mock_okta.go -package=oktaclient --build_flags=--mod=mod

type OktaClient interface {
	CreateUser(ctx context.Context, body okta.CreateUserRequest, qp *query.Params) (*okta.User, *okta.Response, error)
	ListUsers(ctx context.Context, qp *query.Params) ([]*okta.User, *okta.Response, error)
	GetUser(ctx context.Context, userId string) (*okta.User, *okta.Response, error)
	UpdateUser(ctx context.Context, userId string, body okta.User, qp *query.Params) (*okta.User, *okta.Response, error)
	DeactivateOrDeleteUser(ctx context.Context, userId string, qp *query.Params) (*okta.Response, error)
	DeactivateUser(ctx context.Context, userId string, qp *query.Params) (*okta.Response, error)
}

func NewOktaClient(ctx context.Context, cfg *config.OktaConfig) (OktaClient, error) {
	_, client, err := okta.NewClient(
		ctx,
		okta.WithOrgUrl(fmt.Sprintf("https://%s", cfg.Domain)),
		okta.WithToken(cfg.ApiToken),
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to Okta: %s", err.Error())
	}

	return client.User, nil
}
