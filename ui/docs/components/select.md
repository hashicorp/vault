
# Select
Select components are used to render a dropdown.

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| [label] | <code>string</code> | <code>null</code> | The label for the select element. |
| [options] | <code>Array</code> | <code></code> | A list of items that the user will select from. This can be an array of strings or objects. |
| [selectedValue] | <code>string</code> | <code>null</code> | The currently selected item. Can also be used to set the default selected item. This should correspond to the `value` of one of the `<option>`s. |
| [name] | <code>string</code> | <code>null</code> | The name of the select, used for the test selector. |
| [valueAttribute] | <code>string</code> | <code>&quot;value&quot;</code> | When `options` is an array objects, the key to check for when assigning the option elements value. |
| [labelAttribute] | <code>string</code> | <code>&quot;label&quot;</code> | When `options` is an array objects, the key to check for when assigning the option elements' inner text. |
| [isInline] | <code>boolean</code> | <code>false</code> | Whether or not the select should be displayed as inline-block or block. |
| [isFullwidth] | <code>boolean</code> | <code>false</code> | Whether or not the select should take up the full width of the parent element. |
| [noDefault] | <code>boolean</code> | <code>false</code> | shows Select One with empty value as first option |
| [onChange] | <code>Func</code> |  | The action to take once the user has selected an item. This method will be passed the `value` of the select. |

**Example**  
```hbs preview-template
<Select @label='Date Range' @options={{[{ value: 'berry', label: 'Berry' }]}} @onChange={{onChange}}/>
```
