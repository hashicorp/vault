<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/switch.js. To make changes, first edit that file and run "yarn gen-story-md switch" to re-generate the content.-->

## Switch
Switch components are used to indicate boolean values which can be toggled on or off.
They are a stylistic alternative to checkboxes, but use the input[type=checkbox] under the hood.

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| onChange | <code>function</code> |  | onChange is triggered on checkbox change (select, deselect) |
| inputId | <code>string</code> |  | Input ID is needed to match the label to the input |
| [disabled] | <code>boolean</code> | <code>false</code> | disabled makes the switch unclickable |
| [isChecked] | <code>boolean</code> | <code>true</code> | isChecked is the checked status of the input, and must be passed and mutated from the parent |
| [round] | <code>boolean</code> | <code>false</code> | default switch is squared off, this param makes it rounded |
| [size] | <code>string</code> | <code>&quot;&#x27;small&#x27;&quot;</code> | Sizing can be small, medium, or large |
| [status] | <code>string</code> | <code>&quot;&#x27;normal&#x27;&quot;</code> | Status can be normal or success, which makes the switch have a blue background when on |

**Example**
  
```js
<Switch @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
```

**See**

- [Uses of Switch](https://github.com/hashicorp/vault/search?l=Handlebars&q=Switch+OR+switch)
- [Switch Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/switch.js)

---
