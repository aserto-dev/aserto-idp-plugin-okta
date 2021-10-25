package srv

import (
	"context"
	"errors"
	"io"

	"github.com/aserto-dev/aserto-idp-plugin-okta/pkg/config"
	"github.com/aserto-dev/aserto-idp-plugin-okta/pkg/oktaclient"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

type OktaPager func(context.Context, *okta.Response, *[]*okta.User) (*okta.Response, error)

type OktaPlugin struct {
	client       oktaclient.OktaClient
	config       *config.OktaConfig
	pager        OktaPager
	ctx          context.Context
	response     *okta.Response
	users        []*okta.User
	finishedRead bool
}

func (s *OktaPlugin) GetConfig() plugin.PluginConfig {
	return &config.OktaConfig{}
}

func (o *OktaPlugin) Open(cfg plugin.PluginConfig, op plugin.OperationType) error {
	config, ok := cfg.(*config.OktaConfig)

	if !ok {
		return errors.New("invalid config")
	}

	err := o.oktaClient(config)

	if err != nil {
		return err
	}

	o.config = config

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
		oktaUsers, resp, err := o.client.ListUsers(o.ctx, nil)

		if err != nil {
			return nil, err
		}

		for _, u := range oktaUsers {

			user := Transform(u)
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

	resp, err := o.pager(o.ctx, o.response, &(o.users))

	if err != nil {
		errs = multierror.Append(errs, err)
		return nil, errs
	}

	for _, u := range o.users {
		user := Transform(u)
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
	_, _, err := o.client.GetUser(o.ctx, user.Id)

	if err != nil {
		u := TransformToOktaUserReq(user)

		_, _, err := o.client.CreateUser(o.ctx, *u, CreateQueryWithStatus(u.Profile))

		if err != nil {
			return err
		}
	} else {
		updatedUser := &okta.User{
			Profile: ConstructOktaProfile(user),
		}

		_, _, err := o.client.UpdateUser(o.ctx, user.Id, *updatedUser, CreateQueryWithStatus(updatedUser.Profile))

		if err != nil {
			return err
		}
	}

	return nil
}

func (o *OktaPlugin) Delete(id string) error {
	_, err := o.client.DeactivateUser(o.ctx, id, nil)

	if err != nil {
		return err
	}

	_, errs := o.client.DeactivateOrDeleteUser(o.ctx, id, nil)

	if errs != nil {
		return errs
	}

	return nil
}

func (o *OktaPlugin) Close() error {
	return nil
}

func NormalPager() OktaPager {
	return func(ctx context.Context, resp *okta.Response, users *[]*okta.User) (*okta.Response, error) {
		return resp.Next(ctx, users)
	}
}

func (o *OktaPlugin) oktaClient(cfg *config.OktaConfig) error {
	if o.client != nil {
		return nil
	}
	client, err := oktaclient.NewOktaClient(context.Background(), cfg)

	if err != nil {
		return err
	}
	o.client = client

	return nil
}
