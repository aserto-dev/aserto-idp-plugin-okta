package transform

import (
	"reflect"
	"testing"

	"github.com/aserto-dev/aserto-idp-plugin-okta/pkg/testutils"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/stretchr/testify/require"
)

func TestConstructOktaProfile(t *testing.T) {
	assert := require.New(t)
	apiUser := testutils.CreateTestAPIUser("1", "First Last", "testemail@test.com", "active", "40772233223")

	oktaProfile := ConstructOktaProfile(apiUser)

	assert.True(reflect.TypeOf(*oktaProfile) == reflect.TypeOf(okta.UserProfile{}), "the returned object should be *okta.Profile")
	assert.Equal("First", (*oktaProfile)["firstName"], "should correctly detect the first name")
	assert.Equal("Last", (*oktaProfile)["lastName"], "should correctly detect the last name")
	assert.Equal("testemail@test.com", (*oktaProfile)["email"], "should correctly populate the email")
	assert.Equal("testemail@test.com", (*oktaProfile)["login"], "should correctly populate the login")
	assert.Equal("40772233223", (*oktaProfile)["mobilePhone"], "should correctly populate the phone number")
}

func TestConstructOktaProfileWithUserHavingOnlyFirstName(t *testing.T) {
	assert := require.New(t)
	apiUserWOLastName := testutils.CreateTestAPIUser("1", "First", "testemail@test.com", "active", "40772233223")

	oktaProfileWOLastName := ConstructOktaProfile(apiUserWOLastName)

	assert.Equal("First", (*oktaProfileWOLastName)["firstName"], "should detect the first name and not be empty")
	assert.Equal(" ", (*oktaProfileWOLastName)["lastName"], "should detect the last name and be a white space")
}

func TestToOktaUserReq(t *testing.T) {
	assert := require.New(t)
	apiUser := testutils.CreateTestAPIUser("1", "First Last", "testemail@test.com", "active", "40772233223")

	oktaUserReq := ToOktaUserReq(apiUser)

	assert.True(reflect.TypeOf(*oktaUserReq) == reflect.TypeOf(okta.CreateUserRequest{}), "the returned object should be *okta.CreateUserRequest")
	assert.Equal("testemail@test.com", (*oktaUserReq.Profile)["email"], "should correctly populate the user profile")
}

func TestFromOktaWithActiveCompleteUser(t *testing.T) {
	assert := require.New(t)
	oktaUser := testutils.CreateTestOktaUser("1", "ACTIVE", "First", "Last", "testemail@test.com", "+40772233223")
	(*oktaUser.Profile)["additional_info"] = "test"

	apiUser := FromOkta(oktaUser)

	assert.Empty(apiUser.Id, "should leave the id empty")
	assert.Equal("First Last", apiUser.DisplayName, "should correctly construct the displayName")
	assert.Equal("active", apiUser.Attributes.Properties.Fields["status"].GetStringValue(), "should add status to attributes")
	assert.Equal("test", apiUser.Attributes.Properties.Fields["additional_info"].GetStringValue(), "should add additional profile info to attributes")
	assert.Equal(3, len(apiUser.Identities), "3 identities should be populated")
	assert.True(apiUser.Identities["1"].Verified, "should add user id as a verified identity")
	assert.True(apiUser.Identities["40772233223"].Verified, "should add the phone number to identities")
}

func TestFromOktaWithIncorrectPhoneNumberUser(t *testing.T) {
	assert := require.New(t)
	oktaUser := testutils.CreateTestOktaUser("1", "ACTIVE", "First", "Last", "testemail@test.com", "0772233223")

	apiUser := FromOkta(oktaUser)

	assert.Equal(2, len(apiUser.Identities), "2 identities should be populated")
	assert.Nil(apiUser.Identities["0772233223"], "should not add the phone number to identities")
}

func TestFromOktaWithDeprovisionedUser(t *testing.T) {
	assert := require.New(t)
	oktaUser := testutils.CreateTestOktaUser("1", "DEPROVISIONED", "First", "Last", "testemail@test.com", "+40772233223")

	apiUser := FromOkta(oktaUser)

	assert.Empty(apiUser.Id, "should leave the id empty")
	assert.Equal("deprovisioned", apiUser.Attributes.Properties.Fields["status"].GetStringValue(), "should add status to attributes")
	assert.True(apiUser.Identities["1"].Verified, "should always add user id as verified")
	assert.False(apiUser.Identities["testemail@test.com"].Verified, "should add user email as unverified")
	assert.True(apiUser.Deleted, "user should be marked as deleted")
}

func TestFromOktaWithUserCustomAttributes(t *testing.T) {
	assert := require.New(t)

	var roles []interface{}
	roles = append(roles, "admin", "plan")
	oktaUser := testutils.CreateTestOktaUserWithCustomAttribute("roles", roles)

	apiUser := FromOkta(oktaUser)

	rolesTranslated := apiUser.Attributes.Properties.Fields["roles"].GetListValue().Values
	assert.Equal(roles[0], rolesTranslated[0].GetStringValue(), "should add custom attributes with proper type")
}

func TestFromOktaWithInvalidTypeUserCustomAttributes(t *testing.T) {
	assert := require.New(t)

	roles := [2]string{"admin", "plan"}
	oktaUser := testutils.CreateTestOktaUserWithCustomAttribute("roles", roles)

	apiUser := FromOkta(oktaUser)

	rolesTranslated := apiUser.Attributes.Properties.Fields["roles"]
	assert.Nil(rolesTranslated)
}
