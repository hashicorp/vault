---
layout: "docs"
page_title: "OIDC Provider Setup - Auth Methods"
description: |-
  OIDC provider configuration quick starts
---

# OIDC Provider Configuration

This page collects high-level setup steps on how to configure an OIDC application for various
providers. For more general usage and operation information, see the
[Vault JWT/OIDC method documentation](https://www.vaultproject.io/docs/auth/jwt.html).

OIDC providers are often highly configurable and you should become familiar with their
recommended settings and best practices. The instructions below are intended only to help you get
started. Corrections and additions may be submitted via the [Vault Github repository](https://github.com/hashicorp/vault).

## Azure Active Directory (AAD)
Reference: [Azure Active Directory v2.0 and the OpenID Connect protocol](https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-protocols-oidc)

1. Register or select an AAD application. Visit Overview page.
1. Configure Redirect URIs ("Web" type). 
1. Record "Application (client) ID".
1. Under "Endpoints", copy the OpenID Connect metadata document URL, omitting the `/well-known...` portion.
1. Switch to Certificates & Secrets. Create a new client secret and record the generated value as
it will not be accessible after you leave the page.

Please note [Azure AD v2.0 endpoints](https://docs.microsoft.com/en-gb/azure/active-directory/develop/azure-ad-endpoint-comparison) are required for [external groups](https://www.vaultproject.io/docs/secrets/identity/index.html#external-vs-internal-groups) to work. Further, the App Registration needs the [Group.Read.All](https://docs.microsoft.com/en-us/graph/permissions-reference#application-permissions-10)  Microsoft Graph API Permission, and `groupMembershipClaims` should be changed from `none` in the [App registration manifest](https://docs.microsoft.com/en-us/azure/active-directory/develop/reference-app-manifest). In the [OIDC Role config](https://www.vaultproject.io/api/auth/jwt/index.html#create-role) the scope `"https://graph.microsoft.com/.default"` should be added to add groups to the jwt token and `groups_claim` should be set to `groups`. Finally Azure AD group can be referenced by using the groups `objectId` as the [group alias name](https://www.vaultproject.io/api/secret/identity/group-alias.html) for the external group.

## Auth0
1. Select Create Application (Regular Web App).
1. Configure Allowed Callback URLs.
1. Copy client ID and secret.
1. If you see Vault errors involving signature, check the application's Advanced > OAuth settings
 and verify that signing algorithm is "RS256".

## Gitlab
1. Visit Settings > Applications.
1. Fill out Name and Redirect URIs.
1. Making sure to select the "openid" scope.
1. Copy client ID and secret.

## Google
Main reference: [Using OAuth 2.0 to Access Google APIs](https://developers.google.com/identity/protocols/OAuth2)

1. Visit the [Google API Console](https://console.developers.google.com).
1. Create or a select a project.
1. Create a new credential via Credentials > Create Credentials > OAuth Client ID.
1. Configure the OAuth Consent Screen. Application Name is required. Save.
1. Select application type: "Web Application".
1. Configure Authorized Redirect URIs.
1. Save client ID and secret.

## Keycloak
1. Select/create a Realm and Client. Visit Settings.
1. Client Protocol: openid-connect
1. Access Type: confidential
1. Standard Flow Enabled: On
1. Configure Valid Redirect URIs.
1. Visit Settings. Select Client ID and Secret and note the generated secret.

## Okta

1. Make sure an Authorization Server has been created.
1. Visit Applications > Add Application (Web).
1. Configure Login redirect URIs. Save.
1. Save client ID and secret.
