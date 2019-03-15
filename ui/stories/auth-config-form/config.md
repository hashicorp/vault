<a name="AuthConfigForm/Config
The `AuthConfigForm/Config` is the base form to configure auth methods.module_"></a>

## AuthConfigForm/Config
The `AuthConfigForm/Config` is the base form to configure auth methods.

**See**

- [Uses of AuthConfigForm/Config](https://github.com/hashicorp/vault/search?l=Handlebars&q=auth-config-form/config)
- [AuthConfigForm/Config Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/auth-config-form/config.js)

**Properties**

| Name | Type | Default | Description |
| --- | --- | --- | --- |
| model | <code>String</code> | <code></code> | The corresponding auth model that is being configured. |

**Example**

```js
{{auth-config-form/config model.model}}
```
