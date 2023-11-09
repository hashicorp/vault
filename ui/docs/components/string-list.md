# StringList

| Param      | Type                  | Default                        | Description                                                         |
| ---------- | --------------------- | ------------------------------ | ------------------------------------------------------------------- |
| label      | <code>string</code>   |                                | Text displayed in the header above all the inputs.                  |
| onChange   | <code>function</code> |                                | Function called when any of the inputs change.                      |
| inputValue | <code>string</code>   |                                | A string or an array of strings.                                    |
| helpText   | <code>string</code>   |                                | Text displayed as a tooltip.                                        |
| [type]     | <code>string</code>   | <code>&quot;array&quot;</code> | Optional type for inputValue.                                       |
| attrName   | <code>string</code>   |                                | We use this to check the type so we can modify the tooltip content. |
| subText    | <code>string</code>   |                                | Text below the label.                                               |

**Example**

```hbs preview-template
<StringList
  @label='Label here'
  @onChange={{log 'onChange called!'}}
  @inputValue='input value'
  @subText='very important subtext'
/>
```
