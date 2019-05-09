<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/form-field-groups.js. To make changes, first edit that file and run "yarn gen-story-md form-field-groups" to re-generate the content.-->

## FormFieldGroups
`FormFieldGroups` components are field groups associated with a particular model. They render individual `FormField` components.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| [renderGroup] | <code>String</code> | <code></code> | A whitelist of groups to include in the render. |
| model | <code>DS.Model</code> | <code></code> | Model to be passed down to form-field component. If `fieldGroups` is present on the model then it will be iterated over and groups of `FormField` components will be rendered. |
| onChange | <code>Func</code> | <code></code> | Handler that will get set on the `FormField` component. |

**Example**
  
```js
{{if model.fieldGroups}}
 <FormFieldGroups @model={{model}} />
{{/if}}

...

<FormFieldGroups
 @model={{mountModel}}
 @onChange={{action "onTypeChange"}}
 @renderGroup="Method Options"
/>
```

**See**

- [Uses of FormFieldGroups](https://github.com/hashicorp/vault/search?l=Handlebars&q=FormFieldGroups)
- [FormFieldGroups Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/form-field-groups.js)

---
