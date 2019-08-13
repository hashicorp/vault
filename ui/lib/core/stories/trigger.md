<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/confirm.js. To make changes, first edit that file and run "yarn gen-story-md confirm" to re-generate the content.-->

## Trigger
`Trigger` components are a button that shows a confirmation message. They should only be rendered within a `Confirm` component.

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| counters | <code>Array</code> | <code></code> | A list of objects containing the total number of HTTP Requests for each month. `counters` should be the response from the `/internal/counters/requests` endpoint. |

**Example**
  
```js
<div class="box">
  <Confirm as |c|>
    <c.Trigger
      @id={{item.id}}
      @onTrigger={{action c.onTrigger item.id}}
      @triggerText="Delete"
      @message="This will permanently delete this secret and all its vesions."
      @onConfirm={{action "delete" item "secret"}}
      />
  </Confirm>
</div>
```

 * @param id=null {ID} - A unique identifier used to bind a trigger to a confirmation message.
 * @param onTrigger {Func} - A function that displays the confirmation message. This must receive the `id` listed above.
 * @param onConfirm=null {Func} - The action to take when the user clicks the confirm button.
 * @param [title='Delete this?'] {String} - The header text to display in the confirmation message.
 * @param [triggerText='Delete'] {String} - The text on the trigger button.
 * @param [message='You will not be able to recover it later.'] {String} -
 * @param [confirmButtonText='Delete'] {String} - The text to display on the confirm button.
 * @param [cancelButtonText='Cancel'] {String} - The text to display on the cancel button.
 */

**See**

- [Uses of Confirm](https://github.com/hashicorp/vault/search?l=Handlebars&q=Confirm+OR+confirm)
- [Confirm Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/confirm.js)

---
