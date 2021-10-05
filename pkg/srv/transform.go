package srv

const (
	provider = "okta"
)

// func TransformToOkta(in *api.User) (*management.User, error) {
// 	// TODO: add more data here
// 	user := management.User{
// 		ID:           okta.String(in.Id),
// 		Nickname:     okta.String(in.DisplayName),
// 		Email:        okta.String(in.Email),
// 		Picture:      okta.String(in.Picture),
// 		UserMetadata: make(map[string]interface{}),
// 	}
// 	return &user, nil
// }

// // Transform Okta user definition into Aserto Edge User object definition.
// func Transform(in *management.User) (*api.User, error) {

// 	uid := strings.ToLower(strings.TrimPrefix(*in.ID, "okta"))

// 	user := api.User{
// 		Id:          uid,
// 		DisplayName: in.GetNickname(),
// 		Email:       in.GetEmail(),
// 		Picture:     in.GetPicture(),
// 		Identities:  make(map[string]*api.IdentitySource),
// 		Attributes: &api.AttrSet{
// 			Properties:  &structpb.Struct{Fields: make(map[string]*structpb.Value)},
// 			Roles:       []string{},
// 			Permissions: []string{},
// 		},
// 		Applications: make(map[string]*api.AttrSet),
// 		Metadata: &api.Metadata{
// 			CreatedAt: timestamppb.New(in.GetCreatedAt()),
// 			UpdatedAt: timestamppb.New(in.GetUpdatedAt()),
// 		},
// 	}

// 	user.Identities[in.GetID()] = &api.IdentitySource{
// 		Kind:     api.IdentityKind_IDENTITY_KIND_PID,
// 		Provider: provider,
// 		Verified: true,
// 	}

// 	user.Identities[in.GetEmail()] = &api.IdentitySource{
// 		Kind:     api.IdentityKind_IDENTITY_KIND_EMAIL,
// 		Provider: provider,
// 		Verified: in.GetEmailVerified(),
// 	}

// 	phoneProp := strings.ToLower(api.IdentityKind_IDENTITY_KIND_PHONE.String())
// 	if in.UserMetadata[phoneProp] != nil {
// 		phone := in.UserMetadata[phoneProp].(string)
// 		user.Identities[phone] = &api.IdentitySource{
// 			Kind:     api.IdentityKind_IDENTITY_KIND_PHONE,
// 			Verified: false,
// 		}
// 	}

// 	usernameProp := strings.ToLower(api.IdentityKind_IDENTITY_KIND_USERNAME.String())
// 	if in.UserMetadata[usernameProp] != nil {
// 		username := in.UserMetadata[usernameProp].(string)
// 		user.Identities[username] = &api.IdentitySource{
// 			Kind:     api.IdentityKind_IDENTITY_KIND_USERNAME,
// 			Verified: false,
// 		}
// 	}

// 	if in.UserMetadata != nil && len(in.UserMetadata) != 0 {
// 		props, err := structpb.NewStruct(in.UserMetadata)
// 		if err == nil {
// 			user.Attributes.Properties = props
// 		}
// 	}

// 	return &user, nil
// }
