<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/masked-input.js. To make changes, first edit that file and run "yarn gen-story-md masked-input" to re-generate the content.-->

## MaskedInput
`MaskedInput` components are textarea inputs where the input is hidden. They are used to enter sensitive information like passwords.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| [value] | <code>String</code> |  | The value to display in the input. |
| [placeholder] | <code>String</code> | <code>value</code> | The placeholder to display before the user has entered any input. |
| [allowCopy] | <code>bool</code> | <code></code> | Whether or not the input should render with a copy button. |
| [displayOnly] | <code>bool</code> | <code>false</code> | Whether or not to display the value as a display only `pre` element or as an input. |
| [onChange] | <code>function</code> \| <code>action</code> | <code>Function.prototype</code> | A function to call when the value of the input changes. |

**Example**
  
```js
 <MaskedInput
  @value={{attr.options.defaultValue}}
  @placeholder="secret"
  @allowCopy={{true}}
  @onChange={{action "someAction"}}
 />
```

**See**

- [Uses of MaskedInput](https://github.com/hashicorp/vault/search?l=Handlebars&q=MaskedInput)
- [MaskedInput Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/masked-input.js)

---
