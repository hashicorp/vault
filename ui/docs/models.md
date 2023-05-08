# Models

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Decorators and how to use them:](#decorators-and-how-to-use-them)
  - [model-form-fields decorator](#model-form-fields-decorator)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Decorators and how to use them:

### [model-form-fields decorator](../app/decorators/model-form-fields.js)

- Sets `allFields`, `formFields` and/or `formFieldGroups` properties on a model class
- Every model attribute (regardless of args passed to decorator) is expanded and included in the model's `allFields` array.

```js
const formFieldAttrs = ['attrName', 'anotherAttr'];
const formGroupObjects = [
  // Although these keys can be named however you want if using the FormFieldGroups template,
  // default attributes always render and additional keys render inside toggle groups labeled by the key name here
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
  },
```

- the `formFields` array only includes attributes passed to the first argument

```js
model.formFields = [
  {
    name: 'someAttribute',
    type: 'string',
    options: { ...options },
  },
];
```

- the `formFieldGroups` array groups the expanded attributes by key:

```js
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
