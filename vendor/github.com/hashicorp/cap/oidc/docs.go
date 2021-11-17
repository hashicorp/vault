/*
oidc is a package for writing clients that integrate with OIDC Providers using
OIDC flows.


Primary types provided by the package:

* Request: represents one OIDC authentication flow for a user.  It contains the
data needed to uniquely represent that one-time flow across the multiple
interactions needed to complete the OIDC flow the user is attempting.  All
Requests contain an expiration for the user's OIDC flow. Optionally, Requests may
contain overrides of configured provider defaults for audiences, scopes and a
redirect URL.

* Token: represents an OIDC id_token, as well as an Oauth2 access_token and
refresh_token (including the access_token expiry)

* Config: provides the configuration for OIDC provider used by a relying
party (for example: client ID/Secret, redirectURL, supported
signing algorithms, additional scopes requested, etc)

* Provider: provides integration with a provider. The provider provides
capabilities like: generating an auth URL, exchanging codes for tokens,
verifying tokens, making user info requests, etc.

The oidc.callback package

The callback package includes handlers (http.HandlerFunc) which can be used
for the callback leg an OIDC flow. Callback handlers for both the authorization
code flow (with optional PKCE) and the implicit flow are provided.

Example apps

Complete concise example solutions:

* OIDC authentication CLI:
https://github.com/hashicorp/cap/tree/main/oidc/examples/cli/

* OIDC authentication SPA:
https://github.com/hashicorp/cap/tree/main/oidc/examples/spa/

*/
package oidc
