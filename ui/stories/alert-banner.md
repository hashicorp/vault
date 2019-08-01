<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/alert-banner.js. To make changes, first edit that file and run "yarn gen-story-md alert-banner" to re-generate the content.-->

## AlertBanner
`AlertBanner` components are used to inform users of important messages.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| type | <code>String</code> | <code></code> | The banner type. This comes from the message-types helper. |
| [message] | <code>String</code> | <code></code> | The message to display within the banner. |

**Example**
  
```js
<AlertBanner @type="danger" @message="{{model.keyId}} is not a valid lease ID"/>
```

**See**

- [Uses of AlertBanner](https://github.com/hashicorp/vault/search?l=Handlebars&q=AlertBanner)
- [AlertBanner Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/alert-banner.js)

---
