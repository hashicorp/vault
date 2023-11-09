
# Toggle
Toggle components are used to indicate boolean values which can be toggled on or off.
They are a stylistic alternative to checkboxes, but still use the input[type&#x3D;checkbox] under the hood.

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| onChange | <code>function</code> |  | onChange is triggered on checkbox change (select, deselect). Must manually mutate checked value |
| name | <code>string</code> |  | name is passed along to the form field, as well as to generate the ID of the input & "for" value of the label |
| [checked] | <code>boolean</code> | <code>false</code> | checked status of the input, and must be passed in and mutated from the parent |
| [disabled] | <code>boolean</code> | <code>false</code> | disabled makes the switch unclickable |
| [size] | <code>string</code> | <code>&quot;medium&quot;</code> | Sizing can be small or medium |
| [status] | <code>string</code> | <code>&quot;normal&quot;</code> | Status can be normal or success, which makes the switch have a blue background when checked=true |

**Example**  
```hbs preview-template
<Toggle @name="My checked toggle" @checked={{true}}/>
<Toggle @name="Disabled toggle" @disabled={{true}}/>
```
