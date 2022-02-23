package config

import (
	"testing"

	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	"github.com/stretchr/testify/require"
)

func TestValidateWithEmptyDomain(t *testing.T) {
	assert := require.New(t)
	config := OktaConfig{
		Domain:   "",
		APIToken: "token",
	}
	err := config.Validate(plugin.OperationTypeRead)

	assert.NotNil(err)
	assert.Equal("rpc error: code = InvalidArgument desc = no okta domain was provided", err.Error())
}

func TestValidateWithEmptyToken(t *testing.T) {
	assert := require.New(t)
	config := OktaConfig{
		Domain:   "domain",
		APIToken: "",
	}

	err := config.Validate(plugin.OperationTypeRead)

	assert.NotNil(t, err)
	assert.Equal("rpc error: code = InvalidArgument desc = no okta api token was provided", err.Error())
}

func TestValidateWithInvalidCredentials(t *testing.T) {
	assert := require.New(t)
	config := OktaConfig{
		Domain:   "domain",
		APIToken: "token",
	}

	err := config.Validate(plugin.OperationTypeRead)

	assert.NotNil(t, err)
	assert.Contains(err.Error(), "Internal desc = failed to retrieve user from Okta")
}

func TestValidateWithUserPIDAndEmail(t *testing.T) {
	assert := require.New(t)
	config := OktaConfig{
		Domain:    "domain",
		APIToken:  "token",
		UserPID:   "someID",
		UserEmail: "test@email.com",
	}

	err := config.Validate(plugin.OperationTypeWrite)

	assert.NotNil(err)
	assert.Contains(err.Error(), "rpc error: code = InvalidArgument desc = an user PID and an user email were provided; please specify only one")
}

func TestDecription(t *testing.T) {
	assert := require.New(t)
	config := OktaConfig{
		Domain:   "test",
		APIToken: "test",
	}

	description := config.Description()

	assert.Equal("Okta plugin", description, "should return the description of the plugin")
}
