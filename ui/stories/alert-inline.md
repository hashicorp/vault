## AlertInline
`AlertInline` components are used to inform users of important messages.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| [type] | <code>String</code> | <code></code> | The alert type. This comes from the message-types helper. |
| message | <code>String</code> | <code></code> | The message to display within the alert. |

**Example**
  
```js
<AlertInline @type="danger" @message="{{model.keyId}} is not a valid lease ID"/>
```

**See**

- [Uses of AlertInline](https://github.com/hashicorp/vault/search?l=Handlebars&q=AlertInline)
- [AlertInline Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/alert-inline.js) 

---

###### _Documentation generated using [jsdoc-to-markdown](https://github.com/75lb/jsdoc-to-markdown)._
