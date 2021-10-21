package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateWithEmptyDomain(t *testing.T) {
	assert := require.New(t)
	config := OktaConfig{
		OktaDomain:   "",
		OktaApiToken: "test",
	}
	err := config.Validate()

	assert.NotNil(err)
	assert.Equal("rpc error: code = InvalidArgument desc = no okta domain was provided", err.Error())
}

func TestValidateWithEmptyToken(t *testing.T) {
	assert := require.New(t)
	config := OktaConfig{
		OktaDomain:   "dafsdf",
		OktaApiToken: "",
	}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Equal("rpc error: code = InvalidArgument desc = no okta api token was provided", err.Error())
}

func TestDecription(t *testing.T) {
	assert := require.New(t)
	config := OktaConfig{
		OktaDomain:   "test",
		OktaApiToken: "test",
	}

	description := config.Description()

	assert.Equal("Okta plugin", description, "should return the description of the plugin")
}
