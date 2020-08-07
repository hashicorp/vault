<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/alert-banner.js. To make changes, first edit that file and run "yarn gen-story-md alert-banner" to re-generate the content.-->

## AlertBanner
`AlertBanner` components are used to inform users of important messages.

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| type | <code>String</code> | <code></code> | The banner type. This comes from the message-types helper. |
| [secondIconType] | <code>String</code> | <code></code> | If you want a second icon to appear to the right of the title. This comes from the message-types helper. |
| [progressBar] | <code>Object</code> | <code></code> | An object containing a value and maximum for a progress bar. Will be displayed next to the message title. |
| [message] | <code>String</code> | <code></code> | The message to display within the banner. |
| [title] | <code>String</code> | <code></code> | A title to show above the message. If this is not provided, there are default values for each type of alert. |

**Example**
  
```js
<AlertBanner @type="danger" @message="{{model.keyId}} is not a valid lease ID"/>
```

**See**

- [Uses of AlertBanner](https://github.com/hashicorp/vault/search?l=Handlebars&q=AlertBanner+OR+alert-banner)
- [AlertBanner Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/alert-banner.js)

---
