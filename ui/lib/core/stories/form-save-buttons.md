<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/form-save-buttons.js. To make changes, first edit that file and run "yarn gen-story-md form-save-buttons" to re-generate the content.-->

## FormSaveButtons
`FormSaveButtons` displays a button save and a cancel button at the bottom of a form.

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| [saveButtonText] | <code>String</code> | <code>&quot;Save&quot;</code> | The text that will be rendered on the Save button. |
| [isSaving] | <code>Boolean</code> | <code>false</code> | If the form is saving, this should be true. This will disable the save button and render a spinner on it; |
| [cancelLinkParams] | <code>Array</code> | <code>[]</code> | An array of arguments used to construct a link to navigate back to when the Cancel button is clicked. |
| [onCancel] | <code>Fuction</code> | <code></code> | If the form should call an action on cancel instead of route somewhere, the fucntion can be passed using onCancel instead of passing an array to cancelLinkParams. |
| [includeBox] | <code>Boolean</code> | <code>true</code> | By default we include padding around the form with underlines. Passing this value as false will remove that padding. |

**Example**
  
```js
<FormSaveButtons @saveButtonText="Save" @isSaving={{isSaving}} @cancelLinkParams={{array
"foo.route"}} />
```

**See**

- [Uses of FormSaveButtons](https://github.com/hashicorp/vault/search?l=Handlebars&q=FormSaveButtons+OR+form-save-buttons)
- [FormSaveButtons Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/form-save-buttons.js)

---
