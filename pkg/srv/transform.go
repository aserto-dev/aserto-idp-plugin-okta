package srv

import (
	"fmt"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

const (
	provider = "okta"
)

// func TransformToOkta(in *api.User) (*okta.User, error) {
// 	// TODO: add more data here
// 	user := okta.User{
// 		ID:           okta.String(in.Id),
// 		Nickname:     okta.String(in.DisplayName),
// 		Email:        okta.String(in.Email),
// 		Picture:      okta.String(in.Picture),
// 		UserMetadata: make(map[string]interface{}),
// 	}
// 	return &user, nil
// }

// Transform Okta user definition into Aserto Edge User object definition.
func Transform(in *okta.User) (*api.User, error) {

	profileMap := in.Profile
	displayName := fmt.Sprintf("%s %s", (*profileMap)["firstName"], (*profileMap)["lastName"])
	email := fmt.Sprint((*profileMap)["email"])

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
		Verified: true,
	}

	user.Identities[email] = &api.IdentitySource{
		Kind:     api.IdentityKind_IDENTITY_KIND_EMAIL,
		Provider: provider,
		Verified: false,
	}

	if (*profileMap)["mobilePhone"] != nil {
		phone := (*profileMap)["mobilePhone"].(string)
		user.Identities[phone] = &api.IdentitySource{
			Kind:     api.IdentityKind_IDENTITY_KIND_PHONE,
			Verified: false,
		}
	}

	return &user, nil
}
