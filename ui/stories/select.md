<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/select.js. To make changes, first edit that file and run "yarn gen-story-md select" to re-generate the content.-->

## Select
Select components are used to render a dropdown.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| label | <code>String</code> | <code></code> | The label for the select element. |
| [selectedItem] | <code>String</code> | <code></code> | The currently selected item. |
| options | <code>Array</code> | <code></code> | A list of items that the user will select from. |
| [isInline] | <code>Bool</code> | <code>false</code> | Whether or not the select should be displayed as inline-block or block. |
| [isFullwidth] | <code>Bool</code> | <code>false</code> | Whether or not the select should take up the full width of the parent element. |
| onChange | <code>Func</code> | <code></code> | The action to take once the user has selected an item. |

**Example**
  
```js
<Select
  label='Date Range'
  @options={{options}}
  @onChange={{onChange}}/>
```

**See**

- [Uses of Select](https://github.com/hashicorp/vault/search?l=Handlebars&q=Select+OR+select)
- [Select Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/select.js)

---
