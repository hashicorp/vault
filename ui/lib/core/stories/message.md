<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/confirm.js. To make changes, first edit that file and run "yarn gen-story-md confirm" to re-generate the content.-->

## Message
`Message` components trigger and display a confirmation message. They should only be used within a `Confirm` component.

**Properties**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| id | <code>ID</code> | <code></code> | A unique identifier used to bind a trigger to a confirmation message. |
| onConfirm | <code>Func</code> | <code></code> | The action to take when the user clicks the confirm button. |
| [triggerText] | <code>String</code> |<code>'Delete'</code> | The text on the trigger button. |
| [title] | <code>String</code> | <code>'Delete this?'</code> | The header text to display in the confirmation message. |
| [message] | <code>String</code> | <code>'You will not be able to recover it later.'</code> | The message to display above the confirm and cancel buttons. |
| [confirmButtonText] | <code>String</code> | <code>'Delete'</code> | The text to display on the confirm button. |
| [cancelButtonText] | <code>String</code> | <code>'Cancel'</code> | The text to display on the cancel button. |


**Example**
  
```js
<div class="box">
  <Confirm as |c|>
    <c.Message
      @id={{item.id}}
      @triggerText="Delete"
      @message="This will permanently delete this secret and all its versions."
      @onConfirm={{action "delete" item "secret"}}
      />
  </Confirm>
</div>
```

**See**

- [Uses of Confirm](https://github.com/hashicorp/vault/search?l=Handlebars&q=Confirm+OR+confirm)
- [Confirm Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/confirm.js)

---
