package srv

import "github.com/google/wire"

func NewOktaPlugin(cfg *OktaConfig) *OktaPlugin {
	wire.Build(
		wire.Struct(new(OktaPlugin), "*"),
		NewOktaClient,
	)

	return &OktaPlugin{}
}

func NewTestOktaPlugin(cfg *OktaConfig) *OktaPlugin {
	wire.Build(
		wire.Struct(new(OktaPlugin), "*"),
		NewOktaClient,
	)

	return &OktaPlugin{}
}
