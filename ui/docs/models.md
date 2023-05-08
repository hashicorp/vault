# Models

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Decorators and how to use them:](#decorators-and-how-to-use-them)
  - [@withFormFields()](#withformfields)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Decorators and how to use them:

### \* [@withFormFields()](../app/decorators/model-form-fields.js)

- Sets `allFields`, `formFields` and/or `formFieldGroups` properties on a model class
- `allFields` includes every model attribute (regardless of args passed to decorator)
- `formFields` and `formFieldGroups` only exist if the relevant arg is passed to the decorator

```js
const formFieldAttrs = ['attrName', 'anotherAttr'];
const formGroupObjects = [
  // In form-field-groups.hbs form toggle group names are determined by key names
  // 'default' attribute fields render before any toggle groups
  //  additional attribute fields render inside toggle groups
  { default: ['someAttribute'] },
  { 'Additional options': ['anotherAttr'] },
];

@withFormFields(formFieldAttrs, formGroupObjects)
export default class UserModel extends Model {
  @attr('string', { ...options }) someAttribute;
  @attr('boolean', { ...options }) anotherAttr;
}
```

- Each model attribute expands into the following object:

```js
  {
    name: 'someAttribute',
    type: 'string',
    options: { ...options },
  }
```

```js
// only includes attributes passed to the first argument
model.formFields = [
  {
    name: 'someAttribute',
    type: 'string',
    options: { ...options },
  },
];

// expanded attributes are grouped by key
model.formFieldGroups = [
  {
    default: [
      {
        name: 'someAttribute',
        type: 'string',
        options: { ...options },
      },
    ],
  },
  {
    'Additional options': [
      {
        name: 'anotherAttr',
        type: 'boolean',
        options: { ...options },
      },
    ],
  },
];
```
