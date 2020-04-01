<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/ttl-picker2.js. To make changes, first edit that file and run "yarn gen-story-md ttl-picker2" to re-generate the content.-->

## TtlPicker2
TtlPicker2 components are used to enable and select TTL

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| onChange | <code>function</code> |  | This function will be passed a TTL object, which includes enabled{bool}, seconds{number}, timeString{string}. |
| [label] | <code>string</code> | <code>&quot;&#x27;Time to live (TTL)&#x27;&quot;</code> | Label is the main label that lives next to the toggle. |
| [helperTextDisabled] | <code>string</code> | <code>&quot;&#x27;Allow tokens to be used indefinitely&#x27;&quot;</code> | This helper text is shown under the label when the toggle is switched off |
| [helperTextEnabled] | <code>string</code> | <code>&quot;&#x27;Disable the use of the token after&#x27;&quot;</code> | This helper text is shown under the label when the toggle is switched on |
| [time] | <code>number</code> | <code>30</code> | The time (in the default units) which will be adjustable by the user of the form |
| [unit] | <code>string</code> | <code>&quot;&#x27;s&#x27;&quot;</code> | This is the unit key which will show by default on the form. Can be one of `s` (seconds), `m` (minutes), `h` (hours), `d` (days) |

**Example**
  
```js
<TtlPicker2 @onChange={{handleChange}} @time={{defaultTime}} @unit={{defaultUnit}}/>
```

**See**

- [Uses of TtlPicker2](https://github.com/hashicorp/vault/search?l=Handlebars&q=TtlPicker2+OR+ttl-picker2)
- [TtlPicker2 Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/ttl-picker2.js)

---
