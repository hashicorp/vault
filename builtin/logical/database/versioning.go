package database

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/hashicorp/vault/helper/random"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/newdbplugin"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

type databaseVersionWrapper struct {
	database       newdbplugin.Database
	legacyDatabase dbplugin.Database
}

func makeDatabase(ctx context.Context, pluginName string, sys pluginutil.LookRunnerUtil, logger log.Logger) (dbw databaseVersionWrapper, err error) {
	newDB, err := newdbplugin.PluginFactory(ctx, pluginName, sys, logger)
	if err == nil {
		dbw.database = newDB
		return dbw, nil
	}

	legacyDB, err := dbplugin.PluginFactory(ctx, pluginName, sys, logger)
	if err == nil {
		dbw.legacyDatabase = legacyDB
		return dbw, nil
	}

	return dbw, fmt.Errorf("invalid database version")
}

func (db databaseVersionWrapper) Type() (string, error) {
	if db.database != nil {
		return db.database.Type()
	}
	return db.legacyDatabase.Type()
}

func (db databaseVersionWrapper) Close() error {
	if db.database != nil {
		return db.database.Close()
	}
	return db.legacyDatabase.Close()
}

type passwordGenerator interface {
	GeneratePasswordFromPolicy(ctx context.Context, policyName string) (password string, err error)
}

var (
	defaultPasswordGenerator = random.DefaultStringGenerator
)

func generatePassword(ctx context.Context, generator passwordGenerator, passwordPolicy string) (password string, err error) {
	if passwordPolicy == "" {
		return defaultPasswordGenerator.Generate(ctx, rand.Reader)
	}
	return generator.GeneratePasswordFromPolicy(ctx, passwordPolicy)
}

// type userCreator struct {
// 	dbw databaseVersionWrapper
//
// 	createStatements   []string
// 	rollbackStatements []string
// 	displayName        string
// 	roleName           string
// 	expiration         time.Time
//
// 	passwordGenerator passwordGenerator
// 	passwordPolicy    string
// }
//
// func (uc userCreator) createUser(ctx context.Context) (username, password string, err error) {
// 	if uc.dbw.database != nil {
// 		return uc.createUserNew(ctx)
// 	}
//
// }

// func createUser(ctx context.Context, dbw databaseVersionWrapper, createStatements, rollbackStatements []string, displayName, roleName string, expiration time.Time, passwordPolicy string) (username, password string, err error) {
