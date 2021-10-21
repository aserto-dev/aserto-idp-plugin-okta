package srv

import (
	"reflect"
	"testing"

	"github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/stretchr/testify/assert"
)

func TestConstructOktaProfile(t *testing.T) {
	apiUser := CreateTestApiUser("1", "First Last", "testemail@test.com", "active")

	oktaProfile := ConstructOktaProfile(apiUser)
	assert.True(t, reflect.TypeOf(*oktaProfile) == reflect.TypeOf(okta.UserProfile{}), "the returned object should be *okta.Profile")
	assert.Equal(t, (*oktaProfile)["firstName"], "First", "should correctly detect the first name")
	assert.Equal(t, (*oktaProfile)["lastName"], "Last", "should correctly detect the last name")
	assert.Equal(t, (*oktaProfile)["email"], "testemail@test.com", "should correctly populate the email")
	assert.Equal(t, (*oktaProfile)["login"], "testemail@test.com", "should correctly populate the login")
	assert.Equal(t, (*oktaProfile)["status"], "ACTIVE", "should correctly transform and populate the status")

	apiUserWOLastName := CreateTestApiUser("1", "First", "testemail@test.com", "active")
	oktaProfileWOLastName := ConstructOktaProfile(apiUserWOLastName)
	assert.Equal(t, (*oktaProfileWOLastName)["firstName"], "First", "should detect the first name and not be empty")
	assert.Equal(t, (*oktaProfileWOLastName)["lastName"], "", "should detect the last name and be empty")
}

func TestTransformToOktaUserReq(t *testing.T) {
	apiUser := CreateTestApiUser("1", "First Last", "testemail@test.com", "active")

	oktaUserReq := TransformToOktaUserReq(apiUser)

	assert.True(t, reflect.TypeOf(*oktaUserReq) == reflect.TypeOf(okta.CreateUserRequest{}), "the returned object should be *okta.CreateUserRequest")
	assert.Equal(t, (*oktaUserReq.Profile)["email"], "testemail@test.com", "should correctly populate the user profile")
}

func TestTransformWithActiveCompleteUser(t *testing.T) {
	oktaUser := CreateTestOktaUser("1", "active", "First", "Last", "testemail@test.com", "+40772233223")
	apiUser := Transform(oktaUser)

	assert.True(t, reflect.TypeOf(*apiUser) == reflect.TypeOf(api.User{}), "the returned object should be *api.User")
}
