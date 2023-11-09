# MaskedInput

MaskedInput components are textarea inputs where the input is hidden. They are used to enter sensitive information like passwords.

| Param           | Type                                         | Default               | Description                                                                              |
| --------------- | -------------------------------------------- | --------------------- | ---------------------------------------------------------------------------------------- |
| value           | <code>String</code>                          |                       | The value to display in the input.                                                       |
| name            | <code>String</code>                          |                       | The key correlated to the value. Used for the download file name.                        |
| [onChange]      | <code>function</code> \| <code>action</code> | <code>Callback</code> | Callback triggered on change, sends new value. Must set the value of `@value`            |
| [allowCopy]     | <code>boolean</code>                         | <code>false</code>    | Whether or not the input should render with a copy button.                               |
| [allowDownload] | <code>boolean</code>                         | <code>false</code>    | Renders a download button that prompts a confirmation modal to download the secret value |
| [displayOnly]   | <code>boolean</code>                         | <code>false</code>    | Whether or not to display the value as a display only `pre` element or as an input.      |

**Example**

```hbs preview-template
<MaskedInput @value='my secret input' />
<MaskedInput
  @value='some secret display value'
  @displayOnly={{true}}
  @allowCopy={{true}}
  @allowDownload={{true}}
/>
```
