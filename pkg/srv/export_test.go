package srv

import (
	"context"

	"github.com/aserto-dev/aserto-idp-plugin-okta/pkg/oktaclient"
)

// GetClient is ONLY available in test builds,
// making the private client property available to tests.
func (o *OktaPlugin) GetClient() oktaclient.OktaClient {
	return o.client
}

// GetContext is ONLY available in test builds.
// making the private client property available to tests.
func (o *OktaPlugin) GetContext() context.Context {
	return o.ctx
}
