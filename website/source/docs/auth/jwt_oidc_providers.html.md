---
layout: "docs"
page_title: "OIDC Provider Setup - Auth Methods"
description: |-
  OIDC provider configuration quick starts
---

# OIDC Provider Configuration

This page collects high-level setup steps on how to configure an OIDC application
on various providers. Corrections and additions may be submitted via the
[Vault Github repository](https://github.com/hashicorp/vault).

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
1. Configured Authorized Redirect URIs.
1. Save client ID and secret.

## Okta

1. Make sure an Authorization Server has been created.
1. Visit Applications > Add Application (Web).
1. Configure Login redirect URIs. Save.
1. Save client ID and secret.
