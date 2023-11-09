
# ToolbarLink
&#x60;ToolbarLink&#x60; components style links and buttons for the Toolbar
It should only be used inside of &#x60;Toolbar&#x60;.

| Param | Type | Description |
| --- | --- | --- |
| route | <code>string</code> | route to pass to LinkTo |
| model | <code>Model</code> | model to pass to LinkTo |
| models | <code>Array</code> | array of models to pass to LinkTo |
| query | <code>Object</code> | query params to pass to LinkTo |
| replace | <code>boolean</code> | replace arg to pass to LinkTo |
| type | <code>string</code> | Use "add" to change icon to plus sign, or pass in your own kind of icon. |
| disabled | <code>boolean</code> | pass true to disable link |
| disabledTooltip | <code>string</code> | tooltip to display on hover when disabled |

**Example**  
```hbs preview-template
<Toolbar>
  <ToolbarActions>
    <ToolbarLink @route="vault.cluster.policies.create" @type="add" @disabled={{true}} @disabledTooltip="This link is disabled">
      Create policy
    </ToolbarLink>
  </ToolbarActions>
</Toolbar>
```
