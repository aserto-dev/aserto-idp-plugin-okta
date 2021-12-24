package config

import (
	"regexp"
	"testing"

	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	"github.com/stretchr/testify/require"
)

func TestValidateWithEmptyDomain(t *testing.T) {
	assert := require.New(t)
	config := OktaConfig{
		Domain:   "",
		ApiToken: "token",
	}
	err := config.Validate(plugin.OperationTypeRead)

	assert.NotNil(err)
	assert.Equal("rpc error: code = InvalidArgument desc = no okta domain was provided", err.Error())
}

func TestValidateWithEmptyToken(t *testing.T) {
	assert := require.New(t)
	config := OktaConfig{
		Domain:   "domain",
		ApiToken: "",
	}

	err := config.Validate(plugin.OperationTypeRead)

	assert.NotNil(t, err)
	assert.Equal("rpc error: code = InvalidArgument desc = no okta api token was provided", err.Error())
}

func TestValidateWithInvalidCredentials(t *testing.T) {
	assert := require.New(t)
	config := OktaConfig{
		Domain:   "domain",
		ApiToken: "token",
	}

	err := config.Validate(plugin.OperationTypeRead)

	assert.NotNil(t, err)
	r, _ := regexp.Compile("Internal desc = failed to retrieve user from Okta")
	assert.Regexp(r, err.Error())
}

func TestDecription(t *testing.T) {
	assert := require.New(t)
	config := OktaConfig{
		Domain:   "test",
		ApiToken: "test",
	}

	description := config.Description()

	assert.Equal("Okta plugin", description, "should return the description of the plugin")
}
