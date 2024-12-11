// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !minimal

package builtinplugins

import (
	"maps"

	credAliCloud "github.com/hashicorp/vault-plugin-auth-alicloud"
	credAzure "github.com/hashicorp/vault-plugin-auth-azure"
	credCF "github.com/hashicorp/vault-plugin-auth-cf"
	credGcp "github.com/hashicorp/vault-plugin-auth-gcp/plugin"
	credKerb "github.com/hashicorp/vault-plugin-auth-kerberos"
	credKube "github.com/hashicorp/vault-plugin-auth-kubernetes"
	credOCI "github.com/hashicorp/vault-plugin-auth-oci"
	dbCouchbase "github.com/hashicorp/vault-plugin-database-couchbase"
	dbElastic "github.com/hashicorp/vault-plugin-database-elasticsearch"
	dbMongoAtlas "github.com/hashicorp/vault-plugin-database-mongodbatlas"
	dbRedis "github.com/hashicorp/vault-plugin-database-redis"
	dbRedisElastiCache "github.com/hashicorp/vault-plugin-database-redis-elasticache"
	dbSnowflake "github.com/hashicorp/vault-plugin-database-snowflake"
	logicalAd "github.com/hashicorp/vault-plugin-secrets-ad/plugin"
	logicalAlicloud "github.com/hashicorp/vault-plugin-secrets-alicloud"
	logicalAzure "github.com/hashicorp/vault-plugin-secrets-azure"
	logicalGcp "github.com/hashicorp/vault-plugin-secrets-gcp/plugin"
	logicalGcpKms "github.com/hashicorp/vault-plugin-secrets-gcpkms"
	logicalKube "github.com/hashicorp/vault-plugin-secrets-kubernetes"
	logicalMongoAtlas "github.com/hashicorp/vault-plugin-secrets-mongodbatlas"
	logicalLDAP "github.com/hashicorp/vault-plugin-secrets-openldap"
	logicalTerraform "github.com/hashicorp/vault-plugin-secrets-terraform"
	credAws "github.com/hashicorp/vault/builtin/credential/aws"
	credGitHub "github.com/hashicorp/vault/builtin/credential/github"
	credLdap "github.com/hashicorp/vault/builtin/credential/ldap"
	credOkta "github.com/hashicorp/vault/builtin/credential/okta"
	credRadius "github.com/hashicorp/vault/builtin/credential/radius"
	logicalAws "github.com/hashicorp/vault/builtin/logical/aws"
	logicalConsul "github.com/hashicorp/vault/builtin/logical/consul"
	logicalNomad "github.com/hashicorp/vault/builtin/logical/nomad"
	logicalRabbit "github.com/hashicorp/vault/builtin/logical/rabbitmq"
	logicalTotp "github.com/hashicorp/vault/builtin/logical/totp"
	"github.com/hashicorp/vault/helper/pluginconsts"
	dbCass "github.com/hashicorp/vault/plugins/database/cassandra"
	dbHana "github.com/hashicorp/vault/plugins/database/hana"
	dbInflux "github.com/hashicorp/vault/plugins/database/influxdb"
	dbMongo "github.com/hashicorp/vault/plugins/database/mongodb"
	dbMssql "github.com/hashicorp/vault/plugins/database/mssql"
	dbMysql "github.com/hashicorp/vault/plugins/database/mysql"
	dbPostgres "github.com/hashicorp/vault/plugins/database/postgresql"
	dbRedshift "github.com/hashicorp/vault/plugins/database/redshift"
	"github.com/hashicorp/vault/sdk/helper/consts"
)

func newFullAddonRegistry() *registry {
	return &registry{
		credentialBackends: map[string]credentialBackend{
			pluginconsts.AuthTypeAliCloud: {Factory: credAliCloud.Factory},
			pluginconsts.AuthTypeAppId: {
				Factory:           removedFactory,
				DeprecationStatus: consts.Removed,
			},
			pluginconsts.AuthTypeAWS:        {Factory: credAws.Factory},
			pluginconsts.AuthTypeAzure:      {Factory: credAzure.Factory},
			pluginconsts.AuthTypeCF:         {Factory: credCF.Factory},
			pluginconsts.AuthTypeGCP:        {Factory: credGcp.Factory},
			pluginconsts.AuthTypeGitHub:     {Factory: credGitHub.Factory},
			pluginconsts.AuthTypeKerberos:   {Factory: credKerb.Factory},
			pluginconsts.AuthTypeKubernetes: {Factory: credKube.Factory},
			pluginconsts.AuthTypeLDAP:       {Factory: credLdap.Factory},
			pluginconsts.AuthTypeOCI:        {Factory: credOCI.Factory},
			pluginconsts.AuthTypeOkta:       {Factory: credOkta.Factory},
			pluginconsts.AuthTypePCF: {
				Factory:           credCF.Factory,
				DeprecationStatus: consts.Deprecated,
			},
			pluginconsts.AuthTypeRadius: {Factory: credRadius.Factory},
		},
		databasePlugins: map[string]databasePlugin{
			// These four plugins all use the same mysql implementation but with
			// different username settings passed by the constructor.
			"mysql-database-plugin":        {Factory: dbMysql.New(dbMysql.DefaultUserNameTemplate)},
			"mysql-aurora-database-plugin": {Factory: dbMysql.New(dbMysql.DefaultLegacyUserNameTemplate)},
			"mysql-rds-database-plugin":    {Factory: dbMysql.New(dbMysql.DefaultLegacyUserNameTemplate)},
			"mysql-legacy-database-plugin": {Factory: dbMysql.New(dbMysql.DefaultLegacyUserNameTemplate)},

			"cassandra-database-plugin":         {Factory: dbCass.New},
			"couchbase-database-plugin":         {Factory: dbCouchbase.New},
			"elasticsearch-database-plugin":     {Factory: dbElastic.New},
			"hana-database-plugin":              {Factory: dbHana.New},
			"influxdb-database-plugin":          {Factory: dbInflux.New},
			"mongodb-database-plugin":           {Factory: dbMongo.New},
			"mongodbatlas-database-plugin":      {Factory: dbMongoAtlas.New},
			"mssql-database-plugin":             {Factory: dbMssql.New},
			"postgresql-database-plugin":        {Factory: dbPostgres.New},
			"redshift-database-plugin":          {Factory: dbRedshift.New},
			"redis-database-plugin":             {Factory: dbRedis.New},
			"redis-elasticache-database-plugin": {Factory: dbRedisElastiCache.New},
			"snowflake-database-plugin":         {Factory: dbSnowflake.New},
		},
		logicalBackends: map[string]logicalBackend{
			pluginconsts.SecretEngineAD: {
				Factory:           logicalAd.Factory,
				DeprecationStatus: consts.Deprecated,
			},
			pluginconsts.SecretEngineAlicloud: {Factory: logicalAlicloud.Factory},
			pluginconsts.SecretEngineAWS:      {Factory: logicalAws.Factory},
			pluginconsts.SecretEngineAzure:    {Factory: logicalAzure.Factory},
			pluginconsts.SecretEngineCassandra: {
				Factory:           removedFactory,
				DeprecationStatus: consts.Removed,
			},
			pluginconsts.SecretEngineConsul:     {Factory: logicalConsul.Factory},
			pluginconsts.SecretEngineGCP:        {Factory: logicalGcp.Factory},
			pluginconsts.SecretEngineGCPKMS:     {Factory: logicalGcpKms.Factory},
			pluginconsts.SecretEngineKubernetes: {Factory: logicalKube.Factory},
			pluginconsts.SecretEngineMongoDB: {
				Factory:           removedFactory,
				DeprecationStatus: consts.Removed,
			},
			pluginconsts.SecretEngineMongoDBAtlas: {Factory: logicalMongoAtlas.Factory},
			pluginconsts.SecretEngineMSSQL: {
				Factory:           removedFactory,
				DeprecationStatus: consts.Removed,
			},
			pluginconsts.SecretEngineMySQL: {
				Factory:           removedFactory,
				DeprecationStatus: consts.Removed,
			},
			pluginconsts.SecretEngineNomad:    {Factory: logicalNomad.Factory},
			pluginconsts.SecretEngineOpenLDAP: {Factory: logicalLDAP.Factory},
			pluginconsts.SecretEngineLDAP:     {Factory: logicalLDAP.Factory},
			pluginconsts.SecretEnginePostgresql: {
				Factory:           removedFactory,
				DeprecationStatus: consts.Removed,
			},
			pluginconsts.SecretEngineRabbitMQ:  {Factory: logicalRabbit.Factory},
			pluginconsts.SecretEngineTerraform: {Factory: logicalTerraform.Factory},
			pluginconsts.SecretEngineTOTP:      {Factory: logicalTotp.Factory},
		},
	}
}

func extendAddonPlugins(reg *registry) {
	addonReg := newFullAddonRegistry()

	maps.Copy(reg.credentialBackends, addonReg.credentialBackends)
	maps.Copy(reg.databasePlugins, addonReg.databasePlugins)
	maps.Copy(reg.logicalBackends, addonReg.logicalBackends)
}
