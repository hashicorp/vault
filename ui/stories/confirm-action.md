# ConfirmAction

`ConfirmAction` are glyphs used to indicate important information.

## Properties
| Property | Required | Default value | Type | Description | Example |
|---|---|---|---|---|
| onConfirmAction | [x] | `null` | function | The action to take upon confirming. | `() => { console.log('Action!') }` |
| confirmMessage || `'Are you sure you want to do this?'` | string | The message to display upon confirming. ||
| confirmButtonText || `'Delete'` | string | The confirm button text. ||
| cancelButtonText || `'Cancel'` | string | The cancel button text. ||
| disabledMessage || `'Complete the form to complete this action'` | string | The message to display when the button is disabled. ||


## Usage

```javascript
<ConfirmAction
  @onConfirmAction={{ () => { console.log('Action!') } }}
  @confirmMessage="Are you sure you want to delete this config?">
  Delete
</ConfirmAction>
```
https://github.com/hashicorp/vault/search?l=Handlebars&q=ConfirmAction

## Source
https://github.com/hashicorp/vault/blob/master/ui/app/components/confirm-action.js
