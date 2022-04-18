<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/form-field-groups.js. To make changes, first edit that file and run "yarn gen-story-md form-field-groups" to re-generate the content.-->

## onKeyUpCallback : <code>function</code>

**Params**

| Param | Type | Description |
| --- | --- | --- |
| model- | <code>Model</code> | Model to be passed down to form-field component. If `fieldGroups` is present on the model then it will be iterated over and groups of `FormField` components will be rendered. |
| [renderGroup] | <code>string</code> | An allow list of groups to include in the render. |
| [onChange] | <code>onChangeCallback</code> | Handler that will get set on the `FormField` component. |
| [onKeyUp] | [<code>onKeyUpCallback</code>](#onKeyUpCallback) | Handler that will set the value and trigger validation on input changes |
| [modelValidations] | <code>ModelValidations</code> | Object containing validation message for each property |

**Example**
  
```js
{{if model.fieldGroups}}
 <FormFieldGroups @model={{model}} />
{{/if}}

...

<FormFieldGroups
 

**See**

- [Uses of FormFieldGroups](https://github.com/hashicorp/vault/search?l=Handlebars&q=FormFieldGroups+OR+form-field-groups)
- [FormFieldGroups Source Code](https://github.com/hashicorp/vault/blob/main/ui/lib/core/addon/components/form-field-groups.js)

---
