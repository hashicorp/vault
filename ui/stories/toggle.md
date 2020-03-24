<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/toggle.js. To make changes, first edit that file and run "yarn gen-story-md toggle" to re-generate the content.-->

## Toggle
Toggle components are used to indicate boolean values which can be toggled on or off.
They are a stylistic alternative to checkboxes, but still use the input[type=checkbox] under the hood.

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| onChange | <code>function</code> |  | onChange is triggered on checkbox change (select, deselect). Must manually mutate checked value |
| name | <code>string</code> |  | name is passed along to the form field, as well as to generate the ID of the input & "for" value of the label |
| [checked] | <code>boolean</code> | <code>false</code> | checked status of the input, and must be passed in and mutated from the parent |
| [disabled] | <code>boolean</code> | <code>false</code> | disabled makes the switch unclickable |
| [size] | <code>string</code> | <code>&quot;&#x27;medium&#x27;&quot;</code> | Sizing can be small or medium |
| [status] | <code>string</code> | <code>&quot;&#x27;normal&#x27;&quot;</code> | Status can be normal or success, which makes the switch have a blue background when checked=true |

**Example**
  
```js
<Toggle @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
```

**See**

- [Uses of Toggle](https://github.com/hashicorp/vault/search?l=Handlebars&q=Toggle+OR+toggle)
- [Toggle Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/toggle.js)

---
