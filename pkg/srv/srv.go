package srv

import (
	"context"
	"errors"
	"fmt"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// const (
// 	maxBatchSize = int64(500 * 1024)
// )

type OktaPlugin struct {
	Config *OktaConfig
	client *okta.Client
	ctx    context.Context
	// users  []map[string]interface{}
}

func NewOktaPlugin() *OktaPlugin {
	return &OktaPlugin{
		Config: &OktaConfig{},
	}
}

func (s *OktaPlugin) GetConfig() plugin.PluginConfig {
	return &OktaConfig{}
}

func (o *OktaPlugin) Open(cfg plugin.PluginConfig) error {
	config, ok := cfg.(*OktaConfig)
	if !ok {
		return errors.New("invalid config")
	}
	o.Config = config

	ctx, client, err := okta.NewClient(context.Background(),
		okta.WithOrgUrl(fmt.Sprintf("https://%s", config.Domain)),
		okta.WithToken(config.ApiToken),
		okta.WithRequestTimeout(45),
		okta.WithRateLimitMaxRetries(3),
	)

	if err != nil {
		return status.Error(codes.Internal, "failed to connect to Okta")
	}

	o.ctx = ctx
	o.client = client
	return nil
}

func (o *OktaPlugin) Read() ([]*api.User, error) {
	users, _, err := o.client.User.ListUsers(o.ctx, nil)

	for _, user := range users {
		fmt.Println(user)
	}

	return nil, err
}

func (s *OktaPlugin) Write(user *api.User) error {
	// u, err := TransformToOkta(user)
	// if err != nil {
	// 	return err
	// }

	// userMap, size, err := structToMap(u)
	// if err != nil {
	// 	return err
	// }

	// if s.totalSize+size < maxBatchSize {
	// 	s.users = append(s.users, userMap)
	// } else {
	// 	err = s.startJob()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	s.users = make([]map[string]interface{}, 0)
	// 	s.users = append(s.users, userMap)
	// }

	return nil
}

func (s *OktaPlugin) Close() error {
	// if len(s.users) > 0 {
	// 	err := s.startJob()

	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// var errs error
	// for _, j := range s.jobs {
	// 	jobID := okta.StringValue(j.ID)
	// 	err := s.waitJob(jobID)
	// 	if err != nil {
	// 		errs = multierror.Append(errs, err)
	// 	}
	// }
	return nil
}

// func (s *OktaPlugin) waitJob(jobID string) error {
// 	// for {
// 	// 	j, err := s.mgmt.Job.Read(jobID)
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}

// 	// 	switch *j.Status {
// 	// 	case "pending":
// 	// 		{
// 	// 			time.Sleep(1 * time.Second)
// 	// 			continue
// 	// 		}
// 	// 	case "failed":
// 	// 		return fmt.Errorf("Job %s failed", jobID)
// 	// 	case "completed":
// 	// 		return nil
// 	// 	default:
// 	// 		return fmt.Errorf("Unknown status")
// 	// 	}
// 	// }
// 	return nil
// }

// func (s *OktaPlugin) startJob() error {
// 	// job := &management.Job{
// 	// 	ConnectionID:        auth0.String(s.connectionID),
// 	// 	Upsert:              auth0.Bool(true),
// 	// 	SendCompletionEmail: auth0.Bool(false),
// 	// 	Users:               s.users,
// 	// }
// 	// s.wg.Add(1)
// 	// defer s.wg.Done()
// 	// err := s.mgmt.Job.ImportUsers(job)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// s.jobs = append(s.jobs, *job)

// 	return nil
// }

// func structToMap(in interface{}) (map[string]interface{}, int64, error) {
// 	// data, err := json.Marshal(in)
// 	// if err != nil {
// 	// 	return nil, 0, err
// 	// }
// 	// res := make(map[string]interface{})
// 	// err = json.Unmarshal(data, &res)
// 	// if err != nil {
// 	// 	return nil, 0, err
// 	// }
// 	// size := int64(len(data))
// 	//return res, size, nil
// 	return nil, 0, nil
// }
