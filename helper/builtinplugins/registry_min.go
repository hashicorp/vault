// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build minimal

package builtinplugins

import (
	credAliCloud "github.com/hashicorp/vault-plugin-auth-alicloud"
	credCentrify "github.com/hashicorp/vault-plugin-auth-centrify"
	credCF "github.com/hashicorp/vault-plugin-auth-cf"
	credJWT "github.com/hashicorp/vault-plugin-auth-jwt"
	credKerb "github.com/hashicorp/vault-plugin-auth-kerberos"
	credOCI "github.com/hashicorp/vault-plugin-auth-oci"
	dbSnowflake "github.com/hashicorp/vault-plugin-database-snowflake"
	logicalAd "github.com/hashicorp/vault-plugin-secrets-ad/plugin"
	logicalGcp "github.com/hashicorp/vault-plugin-secrets-gcp/plugin"
	logicalGcpKms "github.com/hashicorp/vault-plugin-secrets-gcpkms"
	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	logicalMongoAtlas "github.com/hashicorp/vault-plugin-secrets-mongodbatlas"
	logicalLDAP "github.com/hashicorp/vault-plugin-secrets-openldap"
	credAppRole "github.com/hashicorp/vault/builtin/credential/approle"
	credAws "github.com/hashicorp/vault/builtin/credential/aws"
	credCert "github.com/hashicorp/vault/builtin/credential/cert"
	credGitHub "github.com/hashicorp/vault/builtin/credential/github"
	credLdap "github.com/hashicorp/vault/builtin/credential/ldap"
	credOkta "github.com/hashicorp/vault/builtin/credential/okta"
	credRadius "github.com/hashicorp/vault/builtin/credential/radius"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	logicalAws "github.com/hashicorp/vault/builtin/logical/aws"
	logicalConsul "github.com/hashicorp/vault/builtin/logical/consul"
	logicalNomad "github.com/hashicorp/vault/builtin/logical/nomad"
	logicalPki "github.com/hashicorp/vault/builtin/logical/pki"
	logicalRabbit "github.com/hashicorp/vault/builtin/logical/rabbitmq"
	logicalSsh "github.com/hashicorp/vault/builtin/logical/ssh"
	logicalTotp "github.com/hashicorp/vault/builtin/logical/totp"
	logicalTransit "github.com/hashicorp/vault/builtin/logical/transit"
	dbCass "github.com/hashicorp/vault/plugins/database/cassandra"
	dbHana "github.com/hashicorp/vault/plugins/database/hana"
	dbInflux "github.com/hashicorp/vault/plugins/database/influxdb"
	dbMssql "github.com/hashicorp/vault/plugins/database/mssql"
	dbMysql "github.com/hashicorp/vault/plugins/database/mysql"
	dbPostgres "github.com/hashicorp/vault/plugins/database/postgresql"
	dbRedshift "github.com/hashicorp/vault/plugins/database/redshift"
	"github.com/hashicorp/vault/sdk/helper/consts"
)

func newRegistry() *registry {
	reg := &registry{
		credentialBackends: map[string]credentialBackend{
			"alicloud": {Factory: credAliCloud.Factory},
			"app-id": {
				Factory:           removedFactory,
				DeprecationStatus: consts.Removed,
			},
			"approle": {Factory: credAppRole.Factory},
			"aws":     {Factory: credAws.Factory},
			"centrify": {
				Factory:           credCentrify.Factory,
				DeprecationStatus: consts.Deprecated,
			},
			"cert":     {Factory: credCert.Factory},
			"cf":       {Factory: credCF.Factory},
			"github":   {Factory: credGitHub.Factory},
			"jwt":      {Factory: credJWT.Factory},
			"kerberos": {Factory: credKerb.Factory},
			"ldap":     {Factory: credLdap.Factory},
			"oci":      {Factory: credOCI.Factory},
			"oidc":     {Factory: credJWT.Factory},
			"okta":     {Factory: credOkta.Factory},
			"pcf": {
				Factory:           credCF.Factory,
				DeprecationStatus: consts.Deprecated,
			},
			"radius":   {Factory: credRadius.Factory},
			"userpass": {Factory: credUserpass.Factory},
		},
		databasePlugins: map[string]databasePlugin{
			// These four plugins all use the same mysql implementation but with
			// different username settings passed by the constructor.
			"mysql-database-plugin":        {Factory: dbMysql.New(dbMysql.DefaultUserNameTemplate)},
			"mysql-aurora-database-plugin": {Factory: dbMysql.New(dbMysql.DefaultLegacyUserNameTemplate)},
			"mysql-rds-database-plugin":    {Factory: dbMysql.New(dbMysql.DefaultLegacyUserNameTemplate)},
			"mysql-legacy-database-plugin": {Factory: dbMysql.New(dbMysql.DefaultLegacyUserNameTemplate)},

			"cassandra-database-plugin":  {Factory: dbCass.New},
			"hana-database-plugin":       {Factory: dbHana.New},
			"influxdb-database-plugin":   {Factory: dbInflux.New},
			"mssql-database-plugin":      {Factory: dbMssql.New},
			"postgresql-database-plugin": {Factory: dbPostgres.New},
			"redshift-database-plugin":   {Factory: dbRedshift.New},
			"snowflake-database-plugin":  {Factory: dbSnowflake.New},
		},
		logicalBackends: map[string]logicalBackend{
			"ad": {
				Factory:           logicalAd.Factory,
				DeprecationStatus: consts.Deprecated,
			},
			"aws": {Factory: logicalAws.Factory},
			"cassandra": {
				Factory:           removedFactory,
				DeprecationStatus: consts.Removed,
			},
			"consul": {Factory: logicalConsul.Factory},
			"gcp":    {Factory: logicalGcp.Factory},
			"gcpkms": {Factory: logicalGcpKms.Factory},
			"kv":     {Factory: logicalKv.Factory},
			"mongodb": {
				Factory:           removedFactory,
				DeprecationStatus: consts.Removed,
			},
			// The mongodbatlas secrets engine is not the same as the database plugin equivalent
			// (`mongodbatlas-database-plugin`), and thus will not be deprecated at this time.
			"mongodbatlas": {Factory: logicalMongoAtlas.Factory},
			"mssql": {
				Factory:           removedFactory,
				DeprecationStatus: consts.Removed,
			},
			"mysql": {
				Factory:           removedFactory,
				DeprecationStatus: consts.Removed,
			},
			"nomad": {Factory: logicalNomad.Factory},
			"ldap":  {Factory: logicalLDAP.Factory},
			"pki":   {Factory: logicalPki.Factory},
			"postgresql": {
				Factory:           removedFactory,
				DeprecationStatus: consts.Removed,
			},
			"rabbitmq": {Factory: logicalRabbit.Factory},
			"ssh":      {Factory: logicalSsh.Factory},
			"totp":     {Factory: logicalTotp.Factory},
			"transit":  {Factory: logicalTransit.Factory},
		},
	}

	entAddExtPlugins(reg)

	return reg
}
