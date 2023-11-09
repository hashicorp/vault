# ReadonlyFormField

ReadonlyFormField components are used to display read only, non-editable attributes

| Param | Type                | Description                                                           |
| ----- | ------------------- | --------------------------------------------------------------------- |
| attr  | <code>object</code> | Should be an attribute from a model exported with expandAttributeMeta |
| value | <code>any</code>    | The value that should be displayed on the input                       |

**Example**

```hbs preview-template
<ReadonlyFormField @attr={{hash name='my attr'}} @value='some value' />
```
