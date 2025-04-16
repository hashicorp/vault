package iam

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"github.com/jackc/pgx/v5"
)

func fetchAuthToken(config DBConfig, pgConfig pgx.ConnConfig) (string, error) {
	sess := session.Must(session.NewSession())
	creds := sess.Config.Credentials

	authToken, err := rdsutils.BuildAuthToken(
		fmt.Sprintf("%s:%d", pgConfig.Host, pgConfig.Port),
		config.AWSDBRegion,
		pgConfig.User,
		creds,
	)
	if err != nil {
		return "", err
	}

	config.Logger.Info("new auth token generated for AWS IAM auth", "host", pgConfig.Host, "token", authToken)

	return authToken, nil
}
