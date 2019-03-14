<a name="AuthForm
The `AuthForm` is used to sign users into Vault.module_"></a>

## AuthForm
The `AuthForm` is used to sign users into Vault.

**See**

- [Uses of AuthForm](https://github.com/hashicorp/vault/search?l=Handlebars&q=AuthForm)
- [AuthForm Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/auth-button.js)

**Properties**

| Name | Type | Default | Description |
| --- | --- | --- | --- |
| wrappedToken | <code>String</code> | <code></code> | A token that has been wrapped. |
| cluster | <code>Object</code> | <code></code> | The auth method that is currently selected in the dropdown. This corresponds to an Ember Model. |
| namespace | <code>String</code> | <code></code> | The currently active namespace. |
| redirectTo | <code>String</code> | <code></code> | The name of the route to redirect to. |
| selectedAuth | <code>String</code> | <code></code> | The auth method that is currently selected in the dropdown. |

**Example**

```js
// All properties are passed in via query params.
  <AuthForm 
    @wrappedToken={{wrappedToken}} 
    @cluster={{model}} 
    @namespace={{namespaceQueryParam}} 
    @redirectTo={{redirectTo}} 
    @selectedAuth={{authMethod}}/>```
