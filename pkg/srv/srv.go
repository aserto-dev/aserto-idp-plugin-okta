package srv

import (
	"context"
	"errors"
	"fmt"
	"io"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OktaPlugin struct {
	Config       *OktaConfig
	client       *okta.Client
	ctx          context.Context
	response     *okta.Response
	users        []*okta.User
	finishedRead bool
}

func NewOktaPlugin() *OktaPlugin {
	return &OktaPlugin{
		Config: &OktaConfig{},
	}
}

func (s *OktaPlugin) GetConfig() plugin.PluginConfig {
	return &OktaConfig{}
}

func (o *OktaPlugin) Open(cfg plugin.PluginConfig, op plugin.OperationType) error {
	config, ok := cfg.(*OktaConfig)

	if !ok {
		return errors.New("invalid config")
	}

	o.Config = config

	ctx, client, err := okta.NewClient(context.Background(),
		okta.WithOrgUrl(fmt.Sprintf("https://%s", config.OktaDomain)),
		okta.WithToken(config.OktaApiToken),
	)

	if err != nil {
		return status.Error(codes.Internal, "failed to connect to Okta")
	}

	o.ctx = ctx
	o.client = client
	o.finishedRead = false
	return nil
}

func (o *OktaPlugin) Read() ([]*api.User, error) {
	if o.finishedRead {
		return nil, io.EOF
	}

	var errs error
	var users []*api.User

	if o.response == nil {
		oktaUsers, resp, err := o.client.User.ListUsers(o.ctx, nil)

		if err != nil {
			return nil, err
		}

		for _, u := range oktaUsers {

			user, err := Transform(u)
			if err != nil {
				errs = multierror.Append(errs, err)
			}

			users = append(users, user)
		}

		if resp.HasNextPage() {
			o.response = resp
			o.users = oktaUsers
		} else {
			o.finishedRead = true
		}

		return users, errs
	}

	resp, err := o.response.Next(o.ctx, o.users)

	if err != nil {
		errs = multierror.Append(errs, err)
		return nil, errs
	}

	for _, u := range o.users {
		user, err := Transform(u)
		if err != nil {
			errs = multierror.Append(errs, err)
		}

		users = append(users, user)
	}

	if resp.HasNextPage() {
		o.response = resp
	} else {
		o.finishedRead = true
	}

	return users, errs
}

func (o *OktaPlugin) Write(user *api.User) error {
	_, _, err := o.client.User.GetUser(o.ctx, user.Id)
	qp := query.NewQueryParams(query.WithActivate(true))

	if err != nil {
		u := TransformToOktaUserReq(user)

		_, _, err := o.client.User.CreateUser(o.ctx, *u, qp)

		if err != nil {
			return err
		}
	} else {
		updatedUser := &okta.User{
			Profile: ConstructOktaProfile(user),
		}

		o.client.User.UpdateUser(o.ctx, user.Id, *updatedUser, qp)
	}

	return nil
}

func (s *OktaPlugin) Delete(id string) error {
	return nil
}

func (s *OktaPlugin) Close() error {
	return nil
}
