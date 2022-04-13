<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/form-field.js. To make changes, first edit that file and run "yarn gen-story-md form-field" to re-generate the content.-->

## onKeyUpCallback : <code>function</code>

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| attr | <code>Object</code> |  | usually derived from ember model `attributes` lookup, and all members of `attr.options` are optional |
| model | <code>Model</code> |  | Ember Data model that `attr` is defined on |
| [disabled] | <code>boolean</code> | <code>false</code> | whether the field is disabled |
| [showHelpText] | <code>boolean</code> | <code>true</code> | whether to show the tooltip with help text from OpenAPI |
| [subText] | <code>string</code> |  | text to be displayed below the label |
| [mode] | <code>string</code> |  | used when editType is 'kv' |
| [modelValidations] | <code>ModelValidations</code> |  | Object of errors.  If attr.name is in object and has error message display in AlertInline. |
| [onChange] | <code>onChangeCallback</code> |  | called whenever a value on the model changes via the component |
| [onKeyUp] | [<code>onKeyUpCallback</code>](#onKeyUpCallback) |  | function passed through into MaskedInput to handle validation. It is also handled for certain form-field types here in the action handleKeyUp. |

**Example**
  
```js
{{#each @model.fields as |attr|}}
 <FormField data-test-field @attr={{attr}} @model={{this.model}} />
{{/each}}
```
example attr object:
attr = {
  name: "foo", // name of attribute -- used to populate various fields and pull value from model
  options: {
    label: "Foo", // custom label to be shown, otherwise attr.name will be displayed
    defaultValue: "", // default value to display if model value is not present
    fieldValue: "bar", // used for value lookup on model over attr.name
    editType: "ttl", type of field to use -- example boolean, searchSelect, etc.
    helpText: "This will be in a tooltip",
    readOnly: true
  },
  type: "boolean" // type of attribute value -- string, boolean, etc.
}

**See**

- [Uses of FormField](https://github.com/hashicorp/vault/search?l=Handlebars&q=FormField+OR+form-field)
- [FormField Source Code](https://github.com/hashicorp/vault/blob/main/ui/lib/core/addon/components/form-field.js)

---
