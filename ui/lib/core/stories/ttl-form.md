<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/ttl-form.js. To make changes, first edit that file and run "yarn gen-story-md ttl-form" to re-generate the content.-->

## TtlForm
TtlForm components are used to enter a Time To Live (TTL) input.
This component does not include a label and is designed to take
a time and unit, and pass an object including seconds and
timestring when those two values are changed.

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| onChange | <code>function</code> |  | This function will be called when the user changes the value. An object will be passed in as a parameter with values seconds{number}, timeString{string} |
| [time] | <code>number</code> |  | Time is the value that will be passed into the value input. Can be null/undefined to start if input is required. |
| [unit] | <code>unit</code> | <code>&quot;s&quot;</code> | This is the unit key which will show by default on the form. Can be one of `s` (seconds), `m` (minutes), `h` (hours), `d` (days) |
| [recalculationTimeout] | <code>number</code> | <code>5000</code> | This is the time, in milliseconds, that `recalculateSeconds` will be be true after time is updated |

**Example**
  
```js
<TtlForm @onChange={action handleChange} @unit={{m}}/>
```

**See**

- [Uses of TtlForm](https://github.com/hashicorp/vault/search?l=Handlebars&q=TtlForm+OR+ttl-form)
- [TtlForm Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/ttl-form.js)

---
