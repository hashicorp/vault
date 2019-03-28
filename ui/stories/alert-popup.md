## AlertPopup
The `AlertPopup` is an implementation of the {@link https://github.com/poteto/ember-cli-flash\|ember-cli-flash} `flashMessage`.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| [type] | <code>String</code> | <code></code> | The alert type. This comes from the message-types helper. |
| message | <code>String</code> | <code></code> | The alert message. |
| close | <code>Func</code> | <code></code> | The close action which will close the alert. |

**Example**
  
```js
// All properties are passed in from the flashMessage service.
  <AlertPopup @type={{message-types flash.type}} @message={{flash.message}} @close={{close}}/>```

**See**

- [Uses of AlertPopup](https://github.com/hashicorp/vault/search?l=Handlebars&q=AlertPopup)
- [AlertPopup Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/alert-popup.js) 

---

###### _Documentation generated using [jsdoc-to-markdown](https://github.com/75lb/jsdoc-to-markdown)._
