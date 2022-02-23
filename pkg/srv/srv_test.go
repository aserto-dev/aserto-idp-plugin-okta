package srv

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"testing"

	"github.com/aserto-dev/aserto-idp-plugin-okta/pkg/config"
	"github.com/aserto-dev/aserto-idp-plugin-okta/pkg/oktaclient"
	"github.com/aserto-dev/aserto-idp-plugin-okta/pkg/testutils"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	"github.com/golang/mock/gomock"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/stretchr/testify/require"
)

func TestConstructor(t *testing.T) {
	// Arrange
	assert := require.New(t)

	// Act
	p := NewTestOktaPlugin(gomock.NewController(t), nil)

	// Assert
	assert.NotNil(p)
}

func TestOpenForRead(t *testing.T) {
	// Arrange
	assert := require.New(t)
	p := NewTestOktaPlugin(gomock.NewController(t), nil)

	// Act
	err := p.Open(&config.OktaConfig{}, plugin.OperationTypeRead)

	// Assert
	assert.Nil(err)
}

func TestReadFailToRetriveUserByID(t *testing.T) {
	// Arrange
	assert := require.New(t)
	p := NewTestOktaPlugin(gomock.NewController(t), nil)
	err := p.Open(&config.OktaConfig{UserPID: "invalidID"}, plugin.OperationTypeRead)
	assert.Nil(err)

	p.client.(*oktaclient.MockOktaClient).EXPECT().GetUser(p.ctx, "invalidID").Return(
		nil, nil, errors.New("boom"))

	users, err := p.Read()

	assert.NotNil(err)
	assert.Equal("boom", err.Error(), "should return error")
	assert.Nil(users)
}

func TestReadUserByID(t *testing.T) {
	// Arrange
	assert := require.New(t)
	p := NewTestOktaPlugin(gomock.NewController(t), nil)

	err := p.Open(&config.OktaConfig{UserPID: "userID"}, plugin.OperationTypeRead)
	oktaUser := testutils.CreateTestOktaUser("user1", "active", "stephen", "fry", "stephen@planetexpress.com", "123456")
	assert.Nil(err)

	p.client.(*oktaclient.MockOktaClient).EXPECT().GetUser(p.ctx, "userID").Return(oktaUser, nil, nil)

	users, err := p.Read()

	assert.Nil(err)
	assert.NotNil(users)
	assert.Equal(users[0].Id, oktaUser.Id)
}

func TestReadUserByEmail(t *testing.T) {
	// Arrange
	assert := require.New(t)
	p := NewTestOktaPlugin(gomock.NewController(t), nil)
	err := p.Open(&config.OktaConfig{UserEmail: "stephen@planetexpress.com"}, plugin.OperationTypeRead)
	oktaUser := testutils.CreateTestOktaUser("user1", "active", "stephen", "fry", "stephen@planetexpress.com", "123456")
	assert.Nil(err)

	p.client.(*oktaclient.MockOktaClient).EXPECT().GetUser(p.ctx, "stephen@planetexpress.com").Return(oktaUser, nil, nil)

	users, err := p.Read()

	assert.Nil(err)
	assert.NotNil(users)
	assert.Equal(users[0].Id, oktaUser.Id)
}

func TestReadFailToRetriveUsers(t *testing.T) {
	// Arrange
	assert := require.New(t)
	p := NewTestOktaPlugin(gomock.NewController(t), nil)
	err := p.Open(&config.OktaConfig{}, plugin.OperationTypeRead)
	assert.Nil(err)

	p.client.(*oktaclient.MockOktaClient).EXPECT().ListUsers(p.ctx, gomock.Any()).Return(
		nil, nil, errors.New("boom"))

	users, err := p.Read()

	assert.NotNil(err)
	assert.Equal("boom", err.Error(), "should return error")
	assert.Nil(users)
}

func TestReadSinglePage(t *testing.T) {
	// Arrange
	assert := require.New(t)
	p := NewTestOktaPlugin(gomock.NewController(t), nil)
	err := p.Open(&config.OktaConfig{}, plugin.OperationTypeRead)
	assert.Nil(err)

	p.client.(*oktaclient.MockOktaClient).EXPECT().ListUsers(p.ctx, gomock.Any()).Return([]*okta.User{
		testutils.CreateTestOktaUser("user1", "active", "stephen", "fry", "stephen@planetexpress.com", "123456"),
	}, &okta.Response{NextPage: ""}, nil)

	users, err := p.Read()

	assert.Nil(err)
	assert.Len(users, 1)

	users, err = p.Read()
	assert.NotNil(err)
	assert.Equal(io.EOF, err, "read() should return EOF")
	assert.Nil(users)
}

func TestReadMultiplePagesAneNextPageFail(t *testing.T) {
	// Arrange
	assert := require.New(t)
	p := NewTestOktaPlugin(gomock.NewController(t), func(c context.Context, r *okta.Response, users *[]*okta.User) (*okta.Response, error) {
		return nil, errors.New("boom")
	})
	err := p.Open(&config.OktaConfig{}, plugin.OperationTypeRead)
	assert.Nil(err)

	p.client.(*oktaclient.MockOktaClient).EXPECT().ListUsers(p.ctx, gomock.Any()).Return([]*okta.User{
		testutils.CreateTestOktaUser("user1", "active", "stephen", "fry", "stephen@planetexpress.com", "123456"),
		testutils.CreateTestOktaUser("user2", "active", "stephen2", "fry", "stephen2@planetexpress.com", "123456"),
	}, &okta.Response{NextPage: "yes"}, nil)

	// Act
	users1, err1 := p.Read()
	assert.Nil(err1)
	users2, err2 := p.Read()

	// Assert
	assert.Len(users1, 2)
	assert.Nil(users2)
	assert.NotNil(err2)
	assert.Equal("1 error occurred:\n\t* boom\n\n", err2.Error(), "should return error")
}

func TestReadMultiplePages(t *testing.T) {
	// Arrange
	assert := require.New(t)
	p := NewTestOktaPlugin(gomock.NewController(t), func(c context.Context, r *okta.Response, users *[]*okta.User) (*okta.Response, error) {
		result := []*okta.User{
			testutils.CreateTestOktaUser("user3", "active", "stephen3", "fry", "stephen2@planetexpress.com", "123456"),
		}
		d, _ := json.Marshal(result)
		_ = json.Unmarshal(d, users)
		return &okta.Response{NextPage: ""}, nil
	})
	err := p.Open(&config.OktaConfig{}, plugin.OperationTypeRead)
	assert.Nil(err)

	p.client.(*oktaclient.MockOktaClient).EXPECT().ListUsers(p.ctx, gomock.Any()).Return([]*okta.User{
		testutils.CreateTestOktaUser("user1", "active", "stephen", "fry", "stephen@planetexpress.com", "123456"),
		testutils.CreateTestOktaUser("user2", "active", "stephen2", "fry", "stephen2@planetexpress.com", "123456"),
	}, &okta.Response{NextPage: "yes"}, nil)

	// Act
	users1, err := p.Read()
	assert.Nil(err)
	users2, err := p.Read()

	// Assert
	assert.Nil(err)
	assert.Len(users1, 2)
	assert.Len(users2, 1)
	assert.NotEqual(users1[0].Id, users2[0].Id)
}

func TestDeleteWithInvalidId(t *testing.T) {
	assert := require.New(t)
	p := NewTestOktaPlugin(gomock.NewController(t), nil)
	err := p.Open(&config.OktaConfig{}, plugin.OperationTypeRead)
	assert.Nil(err)

	p.client.(*oktaclient.MockOktaClient).EXPECT().DeactivateUser(p.ctx, "1", nil).Return(nil, errors.New("error"))

	err = p.Delete("1")

	assert.NotNil(err)
	assert.Equal("error", err.Error())
}

func TestDeleteWhenDeleteFails(t *testing.T) {
	assert := require.New(t)
	p := NewTestOktaPlugin(gomock.NewController(t), nil)
	err := p.Open(&config.OktaConfig{}, plugin.OperationTypeRead)
	assert.Nil(err)

	p.client.(*oktaclient.MockOktaClient).EXPECT().DeactivateUser(p.ctx, "1", nil).Return(nil, nil)
	p.client.(*oktaclient.MockOktaClient).EXPECT().DeactivateOrDeleteUser(p.ctx, "1", nil).Return(nil, errors.New("error"))

	err = p.Delete("1")

	assert.NotNil(err)
	assert.Equal("error", err.Error())
}

func TestDeleteSuccess(t *testing.T) {
	assert := require.New(t)
	p := NewTestOktaPlugin(gomock.NewController(t), nil)
	err := p.Open(&config.OktaConfig{}, plugin.OperationTypeRead)
	assert.Nil(err)

	p.client.(*oktaclient.MockOktaClient).EXPECT().DeactivateUser(p.ctx, "1", nil).Return(nil, nil)
	p.client.(*oktaclient.MockOktaClient).EXPECT().DeactivateOrDeleteUser(p.ctx, "1", nil).Return(nil, nil)

	err = p.Delete("1")

	assert.Nil(err)
}

func TestWriteWithNewUserFail(t *testing.T) {
	assert := require.New(t)
	p := NewTestOktaPlugin(gomock.NewController(t), nil)
	err := p.Open(&config.OktaConfig{}, plugin.OperationTypeRead)
	assert.Nil(err)
	user := testutils.CreateTestAPIUser("1", "Name", "mail", "active", "40772233223")

	p.client.(*oktaclient.MockOktaClient).EXPECT().GetUser(p.ctx, user.Id).Return(nil, nil, errors.New("Error1"))
	p.client.(*oktaclient.MockOktaClient).EXPECT().CreateUser(p.ctx, gomock.Any(), gomock.Any()).Return(nil, nil, errors.New("Error2"))

	err = p.Write(user)

	assert.NotNil(err)
	assert.Equal("Error2", err.Error())
}

func TestWriteWithNewUserSuccess(t *testing.T) {
	assert := require.New(t)
	p := NewTestOktaPlugin(gomock.NewController(t), nil)
	err := p.Open(&config.OktaConfig{}, plugin.OperationTypeRead)
	assert.Nil(err)
	user := testutils.CreateTestAPIUser("1", "Name", "mail", "active", "40772233223")

	p.client.(*oktaclient.MockOktaClient).EXPECT().GetUser(p.ctx, user.Id).Return(nil, nil, errors.New("Error1"))
	p.client.(*oktaclient.MockOktaClient).EXPECT().CreateUser(p.ctx, gomock.Any(), gomock.Any()).Return(nil, nil, nil)

	err = p.Write(user)

	assert.Nil(err)
}

func TestWriteWithExistingUserFail(t *testing.T) {
	assert := require.New(t)
	p := NewTestOktaPlugin(gomock.NewController(t), nil)
	err := p.Open(&config.OktaConfig{}, plugin.OperationTypeRead)
	assert.Nil(err)
	user := testutils.CreateTestAPIUser("1", "Name", "mail", "active", "40772233223")

	p.client.(*oktaclient.MockOktaClient).EXPECT().GetUser(p.ctx, user.Id).Return(nil, nil, nil)
	p.client.(*oktaclient.MockOktaClient).EXPECT().UpdateUser(p.ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(
		nil, nil, errors.New("Error"))

	err = p.Write(user)

	assert.NotNil(err)
	assert.Equal("Error", err.Error())
}

func TestWriteWithExistingUserSuccess(t *testing.T) {
	assert := require.New(t)
	p := NewTestOktaPlugin(gomock.NewController(t), nil)
	err := p.Open(&config.OktaConfig{}, plugin.OperationTypeRead)
	assert.Nil(err)
	user := testutils.CreateTestAPIUser("1", "Name", "mail", "active", "40772233223")

	p.client.(*oktaclient.MockOktaClient).EXPECT().GetUser(p.ctx, user.Id).Return(nil, nil, nil)
	p.client.(*oktaclient.MockOktaClient).EXPECT().UpdateUser(p.ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(
		nil, nil, nil)

	err = p.Write(user)

	assert.Nil(err)
}

func TestClose(t *testing.T) {
	assert := require.New(t)
	p := NewTestOktaPlugin(gomock.NewController(t), nil)

	stats, err := p.Close()

	assert.Nil(stats)
	assert.Nil(err)
}
