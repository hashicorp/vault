<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/confirm-action.js. To make changes, first edit that file and run "yarn gen-story-md confirm-action" to re-generate the content.-->

## ConfirmAction
`ConfirmAction` is a button followed by a confirmation message and button used to prevent users from performing actions they do not intend to.

**Properties**

| Name | Type | Default | Description |
| --- | --- | --- | --- |
| onConfirmAction | <code>Func</code> | <code></code> | The action to take upon confirming. |
| [confirmMessage] | <code>String</code> | <code>Are you sure you want to do this?</code> | The message to display upon confirming. |
| [confirmButtonText] | <code>String</code> | <code>Delete</code> | The confirm button text. |
| [cancelButtonText] | <code>String</code> | <code>Cancel</code> | The cancel button text. |
| [disabledMessage] | <code>String</code> | <code>Complete the form to complete this action</code> | The message to display when the button is disabled. |

**Example**
  
```js
<ConfirmAction
  @onConfirmAction={{ () => { console.log('Action!') } }}
  @confirmMessage="Are you sure you want to delete this config?">
  Delete
</ConfirmAction>
 ```
   

**See**

- [Uses of ConfirmAction](https://github.com/hashicorp/vault/search?l=Handlebars&q=ConfirmAction)
- [ConfirmAction Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/confirm-action.js)

---
