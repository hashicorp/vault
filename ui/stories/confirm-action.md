<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/confirm-action.js. To make changes, first edit that file and run "yarn gen-story-md confirm-action" to re-generate the content.-->

## ConfirmAction
`ConfirmAction` is a button followed by a confirmation message and button used to prevent users from performing actions they do not intend to.

**Properties**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| buttonClasses | <code>String</code> | <code></code> | The CSS classes to add to the trigger button. |
| confirmTitle | <code>String</code> | <code>Delete this?</code> | The title to display upon confirming. |
| confirmMessage | <code>String</code> | <code>You will not be able to recover it later.</code> | The message to display upon confirming. |
| confirmButtonText | <code>String</code> | <code>Delete</code> | The confirm button text. |
| cancelButtonText | <code>String</code> | <code>Cancel</code> | The cancel button text. |
| onConfirmAction | <code>Func</code> | <code></code> | The action to take upon confirming. |

**Example**

```js
<ConfirmAction
  @buttonClasses="button is-primary"
  @onConfirmAction={{ () => { console.log('Action!') } }}
  Delete
</ConfirmAction>
 ```


**See**

- [Uses of ConfirmAction](https://github.com/hashicorp/vault/search?l=Handlebars&q=ConfirmAction)
- [ConfirmAction Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/confirm-action.js)

---
