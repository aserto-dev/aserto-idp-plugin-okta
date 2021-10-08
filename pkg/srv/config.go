package srv

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// values set by linker using ldflag -X
var (
	ver    string // nolint:gochecknoglobals // set by linker
	date   string // nolint:gochecknoglobals // set by linker
	commit string // nolint:gochecknoglobals // set by linker
)

func GetVersion() (string, string, string) {
	return ver, date, commit
}

type OktaConfig struct {
	OktaDomain   string `description:"Okta domain" kind:"attribute" mode:"normal" readonly:"false"`
	OktaApiToken string `description:"Okta API Token" kind:"attribute" mode:"normal" readonly:"false"`
}

func (c *OktaConfig) Validate() error {

	if c.OktaDomain == "" {
		return status.Error(codes.InvalidArgument, "no okta domain was provided")
	}

	if c.OktaApiToken == "" {
		return status.Error(codes.InvalidArgument, "no okta api token was provided")
	}

	ctx, client, err := okta.NewClient(context.Background(),
		okta.WithOrgUrl(fmt.Sprintf("https://%s", c.OktaDomain)),
		okta.WithToken(c.OktaApiToken),
	)

	if err != nil {
		return status.Errorf(codes.Internal, "failed to connect to Okta: %s", err.Error())
	}

	filter := query.NewQueryParams(query.WithLimit(1))
	_, _, errReq := client.User.ListUsers(ctx, filter)

	if errReq != nil {
		return status.Errorf(codes.Internal, "failed to retrieve user from Okta: %s", errReq.Error())
	}

	return nil
}

func (c *OktaConfig) Description() string {
	return "Okta plugin"
}
