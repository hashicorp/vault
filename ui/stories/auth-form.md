<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/auth-form.js. To make changes, first edit that file and run "yarn gen-story-md auth-form" to re-generate the content.-->

## AuthForm
The `AuthForm` is used to sign users into Vault.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| wrappedToken | <code>String</code> | <code></code> | The auth method that is currently selected in the dropdown. |
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

**See**

- [Uses of AuthForm](https://github.com/hashicorp/vault/search?l=Handlebars&q=AuthForm)
- [AuthForm Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/auth-form.js)

---
