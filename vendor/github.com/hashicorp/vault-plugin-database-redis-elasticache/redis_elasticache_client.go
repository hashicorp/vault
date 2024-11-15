// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rediselasticache

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/mitchellh/mapstructure"
)

// Verify interface is implemented
var _ dbplugin.Database = (*redisElastiCacheDB)(nil)

type redisElastiCacheDB struct {
	logger hclog.Logger
	config config
	client *elasticache.ElastiCache
}

type config struct {
	AccessKeyID     string `mapstructure:"access_key_id,omitempty"`
	SecretAccessKey string `mapstructure:"secret_access_key,omitempty"`
	Url             string `mapstructure:"url,omitempty"`
	Region          string `mapstructure:"region,omitempty"`

	Username string `mapstructure:"username,omitempty"` // @Deprecated, use AccessKeyID instead
	Password string `mapstructure:"password,omitempty"` // @Deprecated, use SecretAccessKey instead
}

func (r *redisElastiCacheDB) Initialize(_ context.Context, req dbplugin.InitializeRequest) (dbplugin.InitializeResponse, error) {
	r.logger.Debug("initializing AWS ElastiCache Redis client")

	if err := mapstructure.WeakDecode(req.Config, &r.config); err != nil {
		return dbplugin.InitializeResponse{}, err
	}

	// If primary connection attributes are not set, try to fall back on the deprecated values for backward compatibility
	accessKey := r.config.AccessKeyID
	if accessKey == "" && r.config.Username != "" {
		accessKey = r.config.Username
	}
	secretKey := r.config.SecretAccessKey
	if secretKey == "" && r.config.Password != "" {
		secretKey = r.config.Password
	}

	creds, err := awsutil.RetrieveCreds(accessKey, secretKey, "", r.logger)
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("unable to retrieve AWS credentials from provider chain: %w", err)
	}

	region, err := awsutil.GetRegion(r.config.Region)
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("unable to determine AWS region from config nor context: %w", err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	})
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("unable to initialize AWS session: %w", err)
	}
	r.client = elasticache.New(sess)

	if req.VerifyConnection {
		r.logger.Debug("Verifying connection to instance", "url", r.config.Url)

		_, err := r.client.DescribeUsers(nil)
		if err != nil {
			return dbplugin.InitializeResponse{}, fmt.Errorf("unable to connect to ElastiCache Redis endpoint: %w", err)
		}
	}

	return dbplugin.InitializeResponse{
		Config: req.Config,
	}, nil
}

func (r *redisElastiCacheDB) Type() (string, error) {
	return "redisElastiCache", nil
}

func (r *redisElastiCacheDB) Close() error {
	return nil
}

func (r *redisElastiCacheDB) NewUser(_ context.Context, _ dbplugin.NewUserRequest) (dbplugin.NewUserResponse, error) {
	return dbplugin.NewUserResponse{}, fmt.Errorf("user creation not supported")
}

func (r *redisElastiCacheDB) DeleteUser(_ context.Context, _ dbplugin.DeleteUserRequest) (dbplugin.DeleteUserResponse, error) {
	return dbplugin.DeleteUserResponse{}, fmt.Errorf("user deletion not supported")
}

func (r *redisElastiCacheDB) UpdateUser(_ context.Context, req dbplugin.UpdateUserRequest) (dbplugin.UpdateUserResponse, error) {
	r.logger.Debug("updating AWS ElastiCache Redis user", "username", req.Username)

	out, err := r.client.DescribeUsers(&elasticache.DescribeUsersInput{
		UserId: aws.String(req.Username),
	})
	if err != nil {
		return dbplugin.UpdateUserResponse{}, fmt.Errorf("unable to get user %s: %w", req.Username, err)
	}
	if len(out.Users) == 1 && *out.Users[0].Status != "active" {
		return dbplugin.UpdateUserResponse{}, fmt.Errorf("user %s cannot be updated because it is not in the 'active' state", req.Username)
	}

	_, err = r.client.ModifyUser(&elasticache.ModifyUserInput{
		UserId:    &req.Username,
		Passwords: []*string{&req.Password.NewPassword},
	})
	if err != nil {
		return dbplugin.UpdateUserResponse{}, fmt.Errorf("unable to update user %s: %w", req.Username, err)
	}

	return dbplugin.UpdateUserResponse{}, nil
}
