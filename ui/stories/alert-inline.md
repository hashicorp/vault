<a name="AlertInline
`AlertInline` components are used to inform users of important messages.module_"></a>

## AlertInline
`AlertInline` components are used to inform users of important messages.

**See**

- [Uses of AlertInline](https://github.com/hashicorp/vault/search?l=Handlebars&q=AlertInline)
- [AlertInline Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/alert-inline.js)

**Properties**

| Name | Type | Default | Description |
| --- | --- | --- | --- |
| [AlertInline.type] | <code>String</code> | <code></code> | The alert type. This comes from the message-types helper. |
| message | <code>String</code> | <code></code> | The message to display within the alert. |

**Example**

```js
<AlertInline @type="danger" @message="{{model.keyId}} is not a valid lease ID"/>
```
