package srv

import (
	"time"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreateTestApiUser(id, displayName, email, status string) *api.User {
	user := api.User{
		Id:          id,
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
			CreatedAt: timestamppb.New(time.Now()),
			UpdatedAt: timestamppb.New(time.Now()),
		},
	}

	user.Attributes.Properties.Fields["status"] = structpb.NewStringValue(status)

	return &user
}

func CreateTestOktaProfile(firstName, lastName, email, status, mobilePhone string) *okta.UserProfile {
	profile := okta.UserProfile{}
	profile["firstName"] = firstName
	profile["lastName"] = lastName
	profile["email"] = email
	profile["login"] = email
	profile["status"] = status
	profile["mobilePhone"] = mobilePhone

	return &profile
}

func CreateTestOktaUser(id, status, firstName, lastName, email, mobilePhone string) *okta.User {
	now := time.Now()
	profile := CreateTestOktaProfile(firstName, lastName, email, status, mobilePhone)

	user := okta.User{
		Id:          id,
		Created:     &now,
		LastUpdated: &now,
		Status:      status,
		Profile:     profile,
	}

	return &user
}
