<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/ttl-picker2.js. To make changes, first edit that file and run "yarn gen-story-md ttl-picker2" to re-generate the content.-->

## TtlPicker2

TtlPicker2 components are used to enable and select time to live values. Use this TtlPicker2 instead of TtlPicker if you:

- Want the TTL to be enabled or disabled
- Want to have the time recalculated by default when the unit changes (eg 60s -> 1m)

**Params**

| Param                | Type                  | Default                                                  | Description                                                                                                                      |
| -------------------- | --------------------- | -------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| onChange             | <code>function</code> |                                                          | This function will be passed a TTL object, which includes enabled{bool}, seconds{number}, timeString{string}.                    |
| label                | <code>String</code>   | <code>&quot;Time to live (TTL)"</code>                   | Label is the main label that lives next to the toggle.                                                                           |
| helperTextDisabled   | <code>String</code>   | <code>&quot;Allow tokens to be used indefinitely"</code> | This helper text is shown under the label when the toggle is switched off                                                        |
| helperTextEnabled    | <code>String</code>   | <code>&quot;Disable the use of the token after" </code>  | This helper text is shown under the label when the toggle is switched on                                                         |
| time                 | <code>Number</code>   | <code>30</code>                                          | The time (in the default units) which will be adjustable by the user of the form                                                 |
| unit                 | <code>String</code>   | <code>&quot;s&quot;</code>                               | This is the unit key which will show by default on the form. Can be one of `s` (seconds), `m` (minutes), `h` (hours), `d` (days) |
| recalculationTimeout | <code>Number</code>   | <code>5000</code>                                        | This is the time, in milliseconds, that `recalculateSeconds` will be be true after time is updated                               |

**Example**

```js
<TtlPicker2 @onChange={{handleChange}} @time={{defaultTime}} @unit={{defaultUnit}}/>
```

**See**

- [Uses of TtlPicker2](https://github.com/hashicorp/vault/search?l=Handlebars&q=TtlPicker2+OR+ttl-picker2)
- [TtlPicker2 Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/ttl-picker2.js)

---
