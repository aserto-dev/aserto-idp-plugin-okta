package srv

import (
	"context"
	"fmt"
	"log"

	"github.com/okta/okta-sdk-golang/v2/okta"
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
	Domain   string `description:"Okta domain" kind:"attribute" mode:"normal" readonly:"false"`
	ApiToken string `description:"Okta API Token" kind:"attribute" mode:"normal" readonly:"false"`
}

func (c *OktaConfig) Validate() error {
	domain := fmt.Sprintf("https://%s", c.Domain)
	log.Printf("------------------- Starting validation")
	ctx, client, err := okta.NewClient(context.Background(),
		okta.WithOrgUrl(domain),
		okta.WithToken(c.ApiToken),
		okta.WithRequestTimeout(45),
		okta.WithRateLimitMaxRetries(3),
	)

	log.Printf("-------------------  Client created")

	if err != nil {
		return status.Errorf(codes.Internal, "failed to connect to Okta: %s", err.Error())
	}

	users, resp, err := client.User.ListUsers(ctx, nil)

	if err != nil {
		return status.Errorf(codes.Internal, "failed to connect to Okta: %s", err.Error())
	}

	log.Printf("------------------- Users retrieved")
	log.Printf("users: %v\n", users)
	log.Printf("resp: %v\n", resp)
	log.Printf("-------------------  Exit...")

	return nil
}

func (c *OktaConfig) Description() string {
	return "Okta plugin"
}
