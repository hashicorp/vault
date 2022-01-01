package redis

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/vault/helper/random"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfig(b *backend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "config",

			Fields: map[string]*framework.FieldSchema{
				"address": {
					Type:        framework.TypeString,
					Description: "The address of the Redis server",
					Required:    true,
				},

				"username": {
					Type:        framework.TypeString,
					Description: "The username to connect with",
					Required:    true,
				},

				"password": {
					Type:        framework.TypeString,
					Description: "The password to connect with",
				},

				"rotate": {
					Type:        framework.TypeBool,
					Description: "Whether to rotate the password used by Vault immediately",
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.configWrite,
				logical.ReadOperation:   b.configRead,
			},
		},
		{
			Pattern: "config/rotate",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.rotate,
			},
		},
	}
}

type Config struct {
	Address  string
	Username string
	Password string
}

func getConfig(ctx context.Context, s logical.Storage) (*Config, error) {
	entry, err := s.Get(ctx, "config")
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	if entry == nil {
		return nil, nil
	}

	var conf Config
	if err := entry.DecodeJSON(&conf); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return &conf, nil
}

func (c *Config) Client() (*redis.Client, error) {
	if c == nil {
		return nil, fmt.Errorf("the configuration has not been set")
	}

	return redis.NewClient(&redis.Options{
		Addr:     c.Address,
		Username: c.Username,
		Password: c.Password,
	}), nil
}

func (c *Config) Rotate(ctx context.Context) error {
	client, err := c.Client()
	if err != nil {
		return err
	}

	password, err := random.DefaultStringGenerator.Generate(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to generate password: %s", err)
	}

	args := []interface{}{
		"ACL",
		"SETUSER",
		c.Username,
		"#" + hash(password),
	}

	if c.Password != "" {
		args = append(args, "!"+hash(c.Password))
	}

	if _, err := client.Do(ctx, args...).Result(); err != nil {
		return fmt.Errorf("failed to rotate password: %s", err)
	}

	c.Password = password

	return nil
}

func (c *Config) Save(ctx context.Context, s logical.Storage) error {
	entry, err := logical.StorageEntryJSON("config", c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %s", err)
	}

	if err := s.Put(ctx, entry); err != nil {
		return fmt.Errorf("failed to save config: %s", err)
	}

	return nil
}

func (c *Config) Response() *logical.Response {
	if c == nil {
		return logical.ErrorResponse("No configuration found")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"address":  c.Address,
			"username": c.Username,
		},
	}
}

func (b *backend) configWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	conf := &Config{
		Address:  data.Get("address").(string),
		Username: data.Get("username").(string),
		Password: data.Get("password").(string),
	}

	if conf.Address == "" {
		return logical.ErrorResponse("address must be set"), nil
	}
	if conf.Username == "" {
		return logical.ErrorResponse("username must be set"), nil
	}

	client, _ := conf.Client()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return logical.ErrorResponse("failed to connect to the Redis server: %s", err), nil
	}

	if data.Get("rotate").(bool) {
		if err := conf.Rotate(ctx); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	}

	if err := conf.Save(ctx, req.Storage); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	return conf.Response(), nil
}

func (b *backend) configRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	conf, err := getConfig(ctx, req.Storage)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	return conf.Response(), nil
}

func (b *backend) rotate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	conf, err := getConfig(ctx, req.Storage)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	if err := conf.Rotate(ctx); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	if err := conf.Save(ctx, req.Storage); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	return nil, nil
}

func hash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", sum)
}
