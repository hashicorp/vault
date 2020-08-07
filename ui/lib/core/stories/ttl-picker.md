<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/ttl-picker.js. To make changes, first edit that file and run "yarn gen-story-md ttl-picker" to re-generate the content.-->

## TtlPicker
`TtlPicker` components are used to set the time to live.
This version is being deprecated and replaced by `TtlPicker2` which is an automatic-width version that
automatically recalculates the time value when unit is updated unless time has been changed recently.
Once all instances of TtlPicker are replaced with TtlPicker2, this component will be removed and
TtlPicker2 will be renamed to TtlPicker.

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| labelClass | <code>String</code> | <code>&quot;&quot;</code> | A CSS class to add to the label. |
| labelText | <code>String</code> | <code>&quot;TTL&quot;</code> | The text content of the label associated with the widget. |
| initialValue | <code>Number</code> | <code></code> | The starting value of the TTL; |
| setDefaultValue | <code>Boolean</code> | <code>true</code> | If true, the component will trigger onChange on the initial render, causing a value to be set. |
| onChange | <code>function</code> | <code>Function.prototype</code> | The function to call when the value of the ttl changes. |
| outputSeconds | <code>Boolean</code> | <code>false</code> | If true, the component will trigger onChange with a value converted to seconds instead of a Golang duration string. |

**Example**
  
```js
    <TtlPicker @labelText="Lease" @initialValue={{lease}} @onChange={{action (mut lease)}} />
```

**See**

- [Uses of TtlPicker](https://github.com/hashicorp/vault/search?l=Handlebars&q=TtlPicker+OR+ttl-picker)
- [TtlPicker Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/ttl-picker.js)

---
