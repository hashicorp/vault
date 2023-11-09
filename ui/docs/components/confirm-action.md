
# ConfirmAction
ConfirmAction is a button followed by a pop up confirmation message and button used to prevent users from performing actions they do not intend to.

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| [onConfirmAction] | <code>Func</code> | <code></code> | The action to take upon confirming. |
| [confirmTitle] | <code>String</code> | <code>Delete this?</code> | The title to display when confirming. |
| [confirmMessage] | <code>String</code> | <code>You will not be able to recover it later.</code> | The message to display when confirming. |
| [confirmButtonText] | <code>String</code> | <code>Delete</code> | The confirm button text. |
| [cancelButtonText] | <code>String</code> | <code>Cancel</code> | The cancel button text. |
| [buttonClasses] | <code>String</code> |  | A string to indicate the button class. |
| [horizontalPosition] | <code>String</code> | <code>auto-right</code> | For the position of the dropdown. |
| [verticalPosition] | <code>String</code> | <code>below</code> | For the position of the dropdown. |
| [isRunning] | <code>Boolean</code> | <code>false</code> | If action is still running disable the confirm. |
| [disable] | <code>Boolean</code> | <code>false</code> | To disable the confirm action. |

**Example**  
```hbs preview-template
<ConfirmAction @onConfirmAction={{this.myAction}} @confirmMessage="Are you sure you want to delete this config?">
   Delete
 </ConfirmAction>
```
