package srv

import (
	"reflect"
	"testing"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/stretchr/testify/require"
)

func TestConstructOktaProfile(t *testing.T) {
	assert := require.New(t)
	apiUser := CreateTestApiUser("1", "First Last", "testemail@test.com", "active")

	oktaProfile := ConstructOktaProfile(apiUser)

	assert.True(reflect.TypeOf(*oktaProfile) == reflect.TypeOf(okta.UserProfile{}), "the returned object should be *okta.Profile")
	assert.Equal("First", (*oktaProfile)["firstName"], "should correctly detect the first name")
	assert.Equal("Last", (*oktaProfile)["lastName"], "should correctly detect the last name")
	assert.Equal("testemail@test.com", (*oktaProfile)["email"], "should correctly populate the email")
	assert.Equal("testemail@test.com", (*oktaProfile)["login"], "should correctly populate the login")
	assert.Equal("ACTIVE", (*oktaProfile)["status"], "should correctly transform and populate the status")
}

func TestConstructOktaProfileWithUserHavingOnlyFirstName(t *testing.T) {
	assert := require.New(t)
	apiUserWOLastName := CreateTestApiUser("1", "First", "testemail@test.com", "active")

	oktaProfileWOLastName := ConstructOktaProfile(apiUserWOLastName)

	assert.Equal("First", (*oktaProfileWOLastName)["firstName"], "should detect the first name and not be empty")
	assert.Equal("", (*oktaProfileWOLastName)["lastName"], "should detect the last name and be empty")
}

func TestTransformToOktaUserReq(t *testing.T) {
	assert := require.New(t)
	apiUser := CreateTestApiUser("1", "First Last", "testemail@test.com", "active")

	oktaUserReq := TransformToOktaUserReq(apiUser)

	assert.True(reflect.TypeOf(*oktaUserReq) == reflect.TypeOf(okta.CreateUserRequest{}), "the returned object should be *okta.CreateUserRequest")
	assert.Equal("testemail@test.com", (*oktaUserReq.Profile)["email"], "should correctly populate the user profile")
}

func TestTransformWithActiveCompleteUser(t *testing.T) {
	assert := require.New(t)
	oktaUser := CreateTestOktaUser("1", "ACTIVE", "First", "Last", "testemail@test.com", "+40772233223")
	(*oktaUser.Profile)["additional_info"] = "test"

	apiUser := Transform(oktaUser)

	assert.Equal("1", (*apiUser).Id, "should correctly detect the id")
	assert.Equal("First Last", (*apiUser).DisplayName, "should correctly construct the displayName")
	assert.Equal("active", (*apiUser).Attributes.Properties.Fields["status"].GetStringValue(), "should add status to attributes")
	assert.Equal("test", (*apiUser).Attributes.Properties.Fields["additional_info"].GetStringValue(), "should add additional profile info to attributes")
	assert.Equal(3, len((*apiUser).Identities), "3 identities should be populated")
	assert.True((*apiUser).Identities["1"].Verified, "should add user id as a verified identity")
	assert.True((*apiUser).Identities["40772233223"].Verified, "should add the phone number to identities")
}

func TestTransformWithIncorrectPhoneNumberUser(t *testing.T) {
	assert := require.New(t)
	oktaUser := CreateTestOktaUser("1", "ACTIVE", "First", "Last", "testemail@test.com", "0772233223")

	apiUser := Transform(oktaUser)

	assert.Equal(2, len((*apiUser).Identities), "2 identities should be populated")
	assert.Nil((*apiUser).Identities["0772233223"], "should not add the phone number to identities")
}

func TestTransformWithDeactivatedUser(t *testing.T) {
	assert := require.New(t)
	oktaUser := CreateTestOktaUser("1", "DEACTIVATED", "First", "Last", "testemail@test.com", "+40772233223")

	apiUser := Transform(oktaUser)

	assert.Equal("1", (*apiUser).Id, "should correctly detect the id")
	assert.Equal("deactivated", (*apiUser).Attributes.Properties.Fields["status"].GetStringValue(), "should add status to attributes")
	assert.False((*apiUser).Identities["1"].Verified, "should add user id as an unverified identity")
}