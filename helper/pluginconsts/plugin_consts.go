// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pluginconsts

// These consts live outside the plugin registry files to prevent import cycles.
const (
	AuthTypeAliCloud   = "alicloud"
	AuthTypeAppId      = "app-id"
	AuthTypeAWS        = "aws"
	AuthTypeAzure      = "azure"
	AuthTypeCF         = "cf"
	AuthTypeGCP        = "gcp"
	AuthTypeGitHub     = "github"
	AuthTypeKerberos   = "kerberos"
	AuthTypeKubernetes = "kubernetes"
	AuthTypeLDAP       = "ldap"
	AuthTypeOCI        = "oci"
	AuthTypeOkta       = "okta"
	AuthTypePCF        = "pcf"
	AuthTypeRadius     = "radius"
	AuthTypeToken      = "token"
	AuthTypeCert       = "cert"
	AuthTypeOIDC       = "oidc"
	AuthTypeUserpass   = "userpass"
	AuthTypeSAML       = "saml"
	AuthTypeApprole    = "approle"
	AuthTypeJWT        = "jwt"
)
