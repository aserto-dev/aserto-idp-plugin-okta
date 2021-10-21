package srv

import (
	"context"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type OktaUserResource interface {
	ListUsers(ctx context.Context, qp *query.Params) ([]*okta.User, *okta.Response, error)
}

func NewOktaUserResource(oktaClient okta.Client) OktaUserResource {
	return oktaClient.User
}

type MockOktaUserResource struct {
}

type OktaClient interface {
}

type MockOktaClient struct {
	User *MockOktaUserResource
}
