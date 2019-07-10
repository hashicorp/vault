<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/select-dropdown.js. To make changes, first edit that file and run "yarn gen-story-md select-dropdown" to re-generate the content.-->

## SelectDropdown
SelectDropdown components are used to render a dropdown.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| dropdownLabel | <code>String</code> | <code></code> | The label for the select element. |
| [selectedItem] | <code>String</code> | <code></code> | The currently selected item. |
| options | <code>Array</code> | <code></code> | A list of items that the user will select from. |
| [isInline] | <code>Bool</code> | <code>false</code> | Whether or not the select should be displayed as inline-block or block. |
| onChange | <code>Func</code> | <code></code> | The action to take once the user has selected an item. |

**Example**
  
```js
<SelectDropdown
  dropdownLabel='Date Range'
  @options={{options}}
  @onChange={{onChange}}/>
```

**See**

- [Uses of SelectDropdown](https://github.com/hashicorp/vault/search?l=Handlebars&q=SelectDropdown+OR+select-dropdown)
- [SelectDropdown Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/select-dropdown.js)

---
