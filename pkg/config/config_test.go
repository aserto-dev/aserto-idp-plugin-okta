package config_test

import (
	"testing"

	"github.com/aserto-dev/aserto-idp-plugin-okta/pkg/config"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	"github.com/stretchr/testify/require"
)

func TestValidateWithEmptyDomain(t *testing.T) {
	assert := require.New(t)
	cfg := config.OktaConfig{
		Domain:   "",
		APIToken: "token",
	}
	err := cfg.Validate(plugin.OperationTypeRead)

	assert.NotNil(err)
	assert.Equal("rpc error: code = InvalidArgument desc = no okta domain was provided", err.Error())
}

func TestValidateWithEmptyToken(t *testing.T) {
	assert := require.New(t)
	cfg := config.OktaConfig{
		Domain:   "domain",
		APIToken: "",
	}

	err := cfg.Validate(plugin.OperationTypeRead)

	assert.NotNil(t, err)
	assert.Equal("rpc error: code = InvalidArgument desc = no okta api token was provided", err.Error())
}

func TestValidateWithInvalidCredentials(t *testing.T) {
	assert := require.New(t)
	cfg := config.OktaConfig{
		Domain:   "domain",
		APIToken: "token",
	}

	err := cfg.Validate(plugin.OperationTypeRead)

	assert.NotNil(t, err)
	assert.Contains(err.Error(), "Internal desc = failed to retrieve user from Okta")
}

func TestValidateWithUserPIDAndEmail(t *testing.T) {
	assert := require.New(t)
	cfg := config.OktaConfig{
		Domain:    "domain",
		APIToken:  "token",
		UserPID:   "someID",
		UserEmail: "test@email.com",
	}

	err := cfg.Validate(plugin.OperationTypeWrite)

	assert.NotNil(err)
	assert.Contains(err.Error(), "rpc error: code = InvalidArgument desc = an user PID and an user email were provided; please specify only one")
}

func TestDescription(t *testing.T) {
	assert := require.New(t)
	cfg := config.OktaConfig{
		Domain:   "test",
		APIToken: "test",
	}

	description := cfg.Description()

	assert.Equal("Okta plugin", description, "should return the description of the plugin")
}
