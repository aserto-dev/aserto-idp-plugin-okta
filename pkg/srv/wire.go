//+build wireinject

package srv

import (
	"context"

	"github.com/aserto-dev/aserto-idp-plugin-okta/pkg/oktaclient"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/wire"
)

func NewOktaPlugin() *OktaPlugin {
	wire.Build(
		wire.Struct(new(OktaPlugin), "ctx", "pager"),
		context.Background,
		NormalPager,
	)

	return &OktaPlugin{}
}

func NewTestOktaPlugin(ctrl *gomock.Controller, pager OktaPager) *OktaPlugin {
	wire.Build(
		wire.Struct(new(OktaPlugin), "ctx", "client", "pager"),
		wire.Bind(new(oktaclient.OktaClient), new(*oktaclient.MockOktaClient)),
		oktaclient.NewMockOktaClient,
		context.Background,
	)

	return &OktaPlugin{}
}
