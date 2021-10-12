package srv

import (
	"fmt"
	"strings"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dongri/phonenumber"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

const (
	provider = "okta"
)

func ConstructOktaProfile(in *api.User) *okta.UserProfile {
	names := strings.Split(in.DisplayName, " ")
	firstName := names[0]
	lastName := names[1]

	profile := okta.UserProfile{}
	profile["firstName"] = firstName
	profile["lastName"] = lastName
	profile["email"] = in.Email
	profile["login"] = in.Email

	for key, value := range in.Attributes.Properties.Fields {
		if key == "status" {
			profile["status"] = strings.ToUpper(value.GetStringValue())
		} else {
			profile[key] = value
		}
	}

	return &profile
}

func TransformToOktaUserReq(in *api.User) *okta.CreateUserRequest {

	uc := &okta.UserCredentials{}

	user := &okta.CreateUserRequest{
		Credentials: uc,
		Profile:     ConstructOktaProfile(in),
	}

	return user
}

// Transform Okta user definition into Aserto Edge User object definition.
func Transform(in *okta.User) (*api.User, error) {

	profileMap := in.Profile
	displayName := fmt.Sprintf("%s %s", (*profileMap)["firstName"], (*profileMap)["lastName"])
	email := fmt.Sprint((*profileMap)["email"])
	status := strings.ToLower(in.Status)
	verified := false

	switch status {
	case
		"active",
		"recovery",
		"locked out",
		"password expired":
		verified = true
	}

	user := api.User{
		Id:          in.Id,
		DisplayName: displayName,
		Email:       email,
		Picture:     "",
		Identities:  make(map[string]*api.IdentitySource),
		Attributes: &api.AttrSet{
			Properties:  &structpb.Struct{Fields: make(map[string]*structpb.Value)},
			Roles:       []string{},
			Permissions: []string{},
		},
		Applications: make(map[string]*api.AttrSet),
		Metadata: &api.Metadata{
			CreatedAt: timestamppb.New(*in.Created),
			UpdatedAt: timestamppb.New(*in.LastUpdated),
		},
	}

	user.Attributes.Properties.Fields["status"] = structpb.NewStringValue(status)

	for key, value := range *profileMap {
		stringValue := fmt.Sprint(value)

		switch key {
		case
			"mobilePhone",
			"login",
			"email",
			"firstName",
			"lastName":
			continue
		default:
			if value != nil && stringValue != "" {
				user.Attributes.Properties.Fields[key] = structpb.NewStringValue(stringValue)
			}
		}
	}

	user.Identities[in.Id] = &api.IdentitySource{
		Kind:     api.IdentityKind_IDENTITY_KIND_PID,
		Provider: provider,
		Verified: verified,
	}

	user.Identities[email] = &api.IdentitySource{
		Kind:     api.IdentityKind_IDENTITY_KIND_EMAIL,
		Provider: provider,
		Verified: verified,
	}

	mobilePhone := fmt.Sprint((*profileMap)["mobilePhone"])

	country := phonenumber.GetISO3166ByNumber(mobilePhone, true)
	mobilePhoneE164 := phonenumber.ParseWithLandLine(fmt.Sprint((*profileMap)["mobilePhone"]), country.Alpha2)

	if mobilePhoneE164 != "" {
		user.Identities[mobilePhoneE164] = &api.IdentitySource{
			Kind:     api.IdentityKind_IDENTITY_KIND_PHONE,
			Provider: provider,
			Verified: verified,
		}
	}

	return &user, nil
}

func CreateQueryWithStatus(profile *okta.UserProfile) *query.Params {

	if (*profile)["status"] != nil {
		status := fmt.Sprint((*profile)["status"])
		return query.NewQueryParams(query.WithActivate(true), query.WithStatus(status))
	}

	return query.NewQueryParams(query.WithActivate(true))

}
