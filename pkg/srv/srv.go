package srv

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/aserto-dev/aserto-idp-plugin-okta/pkg/config"
	"github.com/aserto-dev/aserto-idp-plugin-okta/pkg/oktaclient"
	"github.com/aserto-dev/aserto-idp-plugin-okta/pkg/transform"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/structpb"
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

func (o *OktaPlugin) GetConfig() plugin.Config {
	return &config.OktaConfig{}
}

func (o *OktaPlugin) GetVersion() (string, string, string) {
	return config.GetVersion()
}

func (o *OktaPlugin) Open(cfg plugin.Config, op plugin.OperationType) error {
	conf, ok := cfg.(*config.OktaConfig)

	if !ok {
		return errors.New("invalid config")
	}

	err := o.oktaClient(conf)

	if err != nil {
		return err
	}

	o.config = conf

	o.finishedRead = false
	return nil
}

func (o *OktaPlugin) Read() ([]*api.User, error) {
	if o.finishedRead {
		return nil, io.EOF
	}

	var errs error
	var users []*api.User

	if o.config.UserPID != "" {
		return o.readByInfo(o.config.UserPID)
	}
	if o.config.UserEmail != "" {
		return o.readByInfo(o.config.UserEmail)
	}

	if o.response == nil {
		oktaUsers, resp, err := o.client.ListUsers(o.ctx, nil)
		if err != nil {
			return nil, err
		}

		if resp.Response != nil && resp.StatusCode == http.StatusTooManyRequests {
			log.Trace().Int("status", resp.StatusCode).Msg("users")
		}

		for _, u := range oktaUsers {
			user := transform.FromOkta(u)

			if err := o.getGroups(u, user); err != nil {
				log.Error().Err(err).Str("userID", u.Id).Msg("getGroups")
			}

			if err := o.getRoles(u, user); err != nil {
				log.Error().Err(err).Str("userID", u.Id).Msg("getRoles")
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

	resp, err := o.pager(o.ctx, o.response, &(o.users))

	if err != nil {
		errs = multierror.Append(errs, err)
		return nil, errs
	}

	for _, u := range o.users {
		user := transform.FromOkta(u)
		users = append(users, user)
	}

	if resp.HasNextPage() {
		o.response = resp
	} else {
		o.finishedRead = true
	}

	return users, errs
}

func (o *OktaPlugin) getGroups(u *okta.User, user *api.User) error {
	if groups, resp, err := o.client.ListUserGroups(o.ctx, u.Id); err == nil && groups != nil && len(groups) != 0 {
		if err != nil {
			return err
		}

		if resp.Response != nil && resp.StatusCode == http.StatusTooManyRequests {
			log.Trace().Int("status", resp.StatusCode).Msg("groups")
		}

		g := make([]interface{}, 0)
		for _, group := range groups {
			g = append(g, group.Profile.Name)
		}

		l, err := structpb.NewList(g)
		if err == nil {
			user.Attributes.Properties.Fields["groups"] = structpb.NewListValue(l)
		}
	}

	return nil
}

func (o *OktaPlugin) getRoles(u *okta.User, user *api.User) error {
	if roles, resp, err := o.client.ListAssignedRolesForUser(o.ctx, u.Id, nil); err == nil && roles != nil && len(roles) != 0 {
		if err != nil {
			return err
		}

		if resp.Response != nil && resp.StatusCode == http.StatusTooManyRequests {
			log.Trace().Int("status", resp.StatusCode).Msg("roles")
		}

		for _, role := range roles {
			user.Attributes.Roles = append(user.Attributes.Roles, role.Type)
		}
	}

	return nil
}

func (o *OktaPlugin) readByInfo(info string) ([]*api.User, error) {
	var users []*api.User

	user, _, err := o.client.GetUser(o.ctx, info)
	o.finishedRead = true
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("failed to get user by %s", info)
	}
	apiUser := transform.FromOkta(user)
	users = append(users, apiUser)
	return users, nil
}

func (o *OktaPlugin) Write(user *api.User) error {
	_, _, err := o.client.GetUser(o.ctx, user.Id)

	if err != nil {
		u := transform.ToOktaUserReq(user)

		_, _, err := o.client.CreateUser(o.ctx, *u, transform.CreateQueryWithStatus(u.Profile))

		if err != nil {
			return err
		}
	} else {
		updatedUser := &okta.User{
			Profile: transform.ConstructOktaProfile(user),
		}

		_, _, err := o.client.UpdateUser(o.ctx, user.Id, *updatedUser, transform.CreateQueryWithStatus(updatedUser.Profile))

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

func (o *OktaPlugin) Close() (*plugin.Stats, error) {
	return nil, nil
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
