
# FormSaveButtons
&#x60;FormSaveButtons&#x60; displays a button save and a cancel button at the bottom of a form.
To show an overall inline error message, use the :error yielded block like shown below.

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| [saveButtonText] | <code>String</code> | <code>&quot;Save&quot;</code> | The text that will be rendered on the Save button. |
| [cancelButtonText] | <code>String</code> | <code>&quot;Cancel&quot;</code> | The text that will be rendered on the Cancel button. |
| [isSaving] | <code>Boolean</code> | <code>false</code> | If the form is saving, this should be true. This will disable the save button and render a spinner on it; |
| [cancelLinkParams] | <code>Array</code> | <code>[]</code> | An array of arguments used to construct a link to navigate back to when the Cancel button is clicked. |
| [onCancel] | <code>function</code> | <code></code> | If the form should call an action on cancel instead of route somewhere, the function can be passed using onCancel instead of passing an array to cancelLinkParams. |
| [includeBox] | <code>Boolean</code> | <code>true</code> | By default we include padding around the form with underlines. Passing this value as false will remove that padding. |

**Example**  
```hbs preview-template
<FormSaveButtons @saveButtonText="Save" @isSaving={{isSaving}} @cancelLinkParams={{array
"foo.route"}}>
  <:error>This is an error</:error>
</FormSaveButtons>
```
