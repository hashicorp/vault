<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/form-field.js. To make changes, first edit that file and run "yarn gen-story-md form-field" to re-generate the content.-->

## FormField
`FormField` components are field elements associated with a particular model.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| [onChange] | <code>Func</code> | <code></code> | Called whenever a value on the model changes via the component. |
| attr | <code>Object</code> | <code></code> | This is usually derived from ember model `attributes` lookup, and all members of `attr.options` are optional. |
| model | <code>DS.Model</code> | <code></code> | The Ember Data model that `attr` is defined on |

### Example Attr

```js
{
   name: "foo",
   options: {
     label: "Foo",
     defaultValue: "",
     editType: "ttl",
     helpText: "This will be in a tooltip"
   },
   type: "boolean"
}
```

**Example**
  
```js
{{#each @model.fields as |attr|}}
  <FormField data-test-field @attr={{attr}} @model={{this.model}} />
{{/each}}
```

**See**

- [Uses of FormField](https://github.com/hashicorp/vault/search?l=Handlebars&q=form-field)
- [FormField Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/form-field.js)

---
