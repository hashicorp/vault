# AuthForm

The `AuthForm` is used to sign users into Vault.

## Properties
| Property | Required | Default value | Type | Description | Example |
|---|---|---|---|---|
| wrappedToken | | `null` | string | The auth method that is currently selected in the dropdown. This is passed in via a query param. | `info` |
| cluster | | `null` | model | The auth method that is currently selected in the dropdown. This is passed in via a query param. | |
| namespace | | `null` | string | The currently active namespace. This is passed in via a query param. | `marketing` |
| redirectTo | | `null` | string | The auth method that is currently selected in the dropdown. This is passed in via a query param. | `info` |
| selectedAuth | | `null` | string | The auth method that is currently selected in the dropdown. This is passed in via a query param. | `info` |

## Usage

```javascript
<AuthForm
  @wrappedToken={{wrappedToken}}
  @cluster={{model}}
  @namespace={{namespaceQueryParam}}
  @redirectTo={{redirectTo}}
  @selectedAuth="userpass"
/>
```
https://github.com/hashicorp/vault/search?l=Handlebars&q=AuthForm

## Source
https://github.com/hashicorp/vault/blob/master/ui/app/components/auth-form.js
s.g7zd4UQlTWP4DyAdRa3PNPjT
