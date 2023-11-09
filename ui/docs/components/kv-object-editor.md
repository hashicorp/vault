# KvObjectEditor

KvObjectEditor components are called in FormFields when the editType on the model is kv. They are used to show a key-value input field.

| Param              | Type                  | Default            | Description                                                                                                                  |
| ------------------ | --------------------- | ------------------ | ---------------------------------------------------------------------------------------------------------------------------- |
| value              | <code>string</code>   |                    | the value is captured from the model.                                                                                        |
| onChange           | <code>function</code> |                    | function that captures the value on change                                                                                   |
| [isMasked]         | <code>boolean</code>  | <code>false</code> | when true the `<MaskedInput>` renders instead of the default `<textarea>` to input the value portion of the key/value object |
| [onKeyUp]          | <code>function</code> |                    | function passed in that handles the dom keyup event. Used for validation on the kv custom metadata.                          |
| [label]            | <code>string</code>   |                    | label displayed over key value inputs                                                                                        |
| [labelClass]       | <code>string</code>   |                    | override default label class in FormFieldLabel component                                                                     |
| [warning]          | <code>string</code>   |                    | warning that is displayed                                                                                                    |
| [helpText]         | <code>string</code>   |                    | helper text. In tooltip.                                                                                                     |
| [subText]          | <code>string</code>   |                    | placed under label.                                                                                                          |
| [keyPlaceholder]   | <code>string</code>   |                    | placeholder for key input                                                                                                    |
| [valuePlaceholder] | <code>string</code>   |                    | placeholder for value input                                                                                                  |

**Example**

```hbs preview-template
<KvObjectEditor @value={{hash foo='bar'}} @onChange={{log 'on change called!'}} @label='Label here' />
```
