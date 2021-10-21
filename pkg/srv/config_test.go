package srv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateWithEmptyCredentials(t *testing.T) {

	config := OktaConfig{
		OktaDomain:   "",
		OktaApiToken: "",
	}
	err := config.Validate()

	assert.NotNil(t, err)

	config = OktaConfig{
		OktaDomain:   "dafsdf",
		OktaApiToken: "",
	}

	err = config.Validate()

	assert.NotNil(t, err)

}
