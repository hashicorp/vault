# FormFieldLabel

FormFieldLabel components add labels and descriptions to inputs

| Param      | Type                | Description                                                    |
| ---------- | ------------------- | -------------------------------------------------------------- |
| [label]    | <code>string</code> | label text -- component attributes are spread on label element |
| [helpText] | <code>string</code> | adds a tooltip                                                 |
| [subText]  | <code>string</code> | de-emphasized text rendered below the label                    |
| [docLink]  | <code>string</code> | url to documentation rendered after the subText                |

**Example**

```hbs preview-template
<FormFieldLabel
  for='input-name'
  @label='Label'
  @helpText='Important help text'
  @subText='Subtext is here'
  @docLink='/vault'
/>
```
