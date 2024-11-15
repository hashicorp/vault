// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rediselasticache

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
)

func New() (interface{}, error) {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	db := &redisElastiCacheDB{
		logger: logger,
	}

	return wrapWithSanitizerMiddleware(db), nil
}

func wrapWithSanitizerMiddleware(db *redisElastiCacheDB) dbplugin.Database {
	return dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValuesToMask)
}

func (r *redisElastiCacheDB) secretValuesToMask() map[string]string {
	return map[string]string{
		r.config.Password:        "[password]",
		r.config.Username:        "[username]",
		r.config.AccessKeyID:     "[access_key_id]",
		r.config.SecretAccessKey: "[secret_access_key]",
	}
}
