<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/auth-form.js. To make changes, first edit that file and run "yarn gen-story-md auth-form" to re-generate the content.-->

## AuthForm
The `AuthForm` is used to sign users into Vault.

**Params**

| Param | Type | Description |
| --- | --- | --- |
| wrappedToken | <code>string</code> | The auth method that is currently selected in the dropdown. |
| cluster | <code>object</code> | The auth method that is currently selected in the dropdown. This corresponds to an Ember Model. |
| namespace- | <code>string</code> | The currently active namespace. |
| selectedAuth | <code>string</code> | The auth method that is currently selected in the dropdown. |
| onSuccess | <code>function</code> | Fired on auth success |

**Example**
  
```js
// All properties are passed in via query params.
<AuthForm @wrappedToken={{wrappedToken}} @cluster={{model}} @namespace={{namespaceQueryParam}} @selectedAuth={{authMethod}} @onSuccess={{action this.onSuccess}} />```

**See**

- [Uses of AuthForm](https://github.com/hashicorp/vault/search?l=Handlebars&q=AuthForm+OR+auth-form)
- [AuthForm Source Code](https://github.com/hashicorp/vault/blob/main/ui/app/components/auth-form.js)

---
