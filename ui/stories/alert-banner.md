<a name="AlertBanner
`AlertBanner` components are used to inform users of important messages.module_"></a>

## AlertBanner
`AlertBanner` components are used to inform users of important messages.

**See**

- [Uses of AlertBanner](https://github.com/hashicorp/vault/search?l=Handlebars&q=AlertBanner)
- [AlertBanner Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/alert-banner.js)

**Properties**

| Name | Type | Default | Description |
| --- | --- | --- | --- |
| [AlertBanner.type] | <code>String</code> | <code></code> | The banner type. This comes from the message-types helper. |
| message | <code>String</code> | <code></code> | The message to display within the banner. |

**Example**

```js
<AlertBanner @type="danger" @message="{{model.keyId}} is not a valid lease ID"/>
```
