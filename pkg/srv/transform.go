package srv

import (
	"fmt"
	"strings"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/okta/okta-sdk-golang/v2/okta"
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

	return &profile
}

func TransformToOktaUserReq(in *api.User) *okta.CreateUserRequest {
	// TODO: add phoneNumber, status & custom attributes to Profile

	uc := &okta.UserCredentials{}

	user := &okta.CreateUserRequest{
		Credentials: uc,
		Profile:     ConstructOktaProfile(in),
	}

	return user
}

// Transform Okta user definition into Aserto Edge User object definition.
func Transform(in *okta.User) (*api.User, error) {
	// TODO: add status as an attribute & other custom attributes

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

	if (*profileMap)["mobilePhone"] != nil {
		phone := (*profileMap)["mobilePhone"].(string)
		user.Identities[phone] = &api.IdentitySource{
			Kind:     api.IdentityKind_IDENTITY_KIND_PHONE,
			Verified: verified,
		}
	}

	return &user, nil
}
