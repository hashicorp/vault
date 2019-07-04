<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/message-error.js. To make changes, first edit that file and run "yarn gen-story-md message-error" to re-generate the content.-->

## MessageError
`MessageError` extracts an error from a model or a passed error and displays it using the `AlertBanner` component.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| model | <code>DS.Model</code> | <code></code> | An Ember data model that will be used to bind error statest to the internal `errors` property. |
| errors | <code>Array</code> | <code></code> | An array of error strings to show. |
| errorMessage | <code>String</code> | <code></code> | An Error string to display. |

**Example**
  
```js
<MessageError @model={{model}} />
```

**See**

- [Uses of MessageError](https://github.com/hashicorp/vault/search?l=Handlebars&q=MessageError)
- [MessageError Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/message-error.js)

---
