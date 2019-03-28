
[//]: # To adjust what shows up in generated markdown files, select partials from https://github.com/jsdoc2md/dmd/tree/master/partials and/or write your own version (as with examples below)

## AlertBanner
`AlertBanner` components are used to inform users of important messages.

**Properties**

| Name | Type | Default | Description |
| --- | --- | --- | --- |
| type | <code>String</code> | <code></code> | The banner type. This comes from the message-types helper. |
| message | <code>String</code> | <code></code> | The message to display within the banner. |

**Example**
  
```js
<AlertBanner @type="danger" @message="{{model.keyId}} is not a valid lease ID"/>
```

**See**

- [Uses of AlertBanner](https://github.com/hashicorp/vault/search?l=Handlebars&q=AlertBanner)
- [AlertBanner Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/alert-banner.js) 

---

###### _Documentation generated using [jsdoc-to-markdown](https://github.com/75lb/jsdoc-to-markdown)._ ######
