<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/alert-inline.js. To make changes, first edit that file and run "yarn gen-story-md alert-inline" to re-generate the content.-->

## AlertInline
`AlertInline` components are used to inform users of important messages.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| type | <code>String</code> | <code></code> | The alert type. This comes from the message-types helper. |
| [message] | <code>String</code> | <code></code> | The message to display within the alert. |

**Example**
  
```js
<AlertInline
  @type="danger"
  @message="{{model.keyId}} is not a valid lease ID"/>
```

**See**

- [Uses of AlertInline](https://github.com/hashicorp/vault/search?l=Handlebars&q=AlertInline)
- [AlertInline Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/alert-inline.js)

---
