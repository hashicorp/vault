# ToolbarLink

ToolbarLink components style links and buttons for the Toolbar
It should only be used inside of Toolbar component.

| Param           | Type                 | Description                                                              |
| --------------- | -------------------- | ------------------------------------------------------------------------ |
| route           | <code>string</code>  | route to pass to LinkTo                                                  |
| model           | <code>model</code>   | model to pass to LinkTo                                                  |
| models          | <code>array</code>   | array of models to pass to LinkTo                                        |
| query           | <code>object</code>  | query params to pass to LinkTo                                           |
| replace         | <code>boolean</code> | replace arg to pass to LinkTo                                            |
| type            | <code>string</code>  | Use "add" to change icon to plus sign, or pass in your own kind of icon. |
| disabled        | <code>boolean</code> | pass true to disable link                                                |
| disabledTooltip | <code>string</code>  | tooltip to display on hover when disabled                                |

**Example**

```hbs preview-template
<Toolbar>
  <ToolbarActions>
    <ToolbarLink @route='vault' @disabled={{true}} @disabledTooltip='This link is disabled'>
      Disabled link
    </ToolbarLink>
    <ToolbarLink @route='vault' @type='add'>
      Create item
    </ToolbarLink>
  </ToolbarActions>
</Toolbar>
```
