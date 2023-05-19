# Models

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Models](#models)
  - [Capabilities](#capabilities)
  - [Decorators](#decorators)
    - [@withFormFields()](#withformfields)
    - [@withModelValidations()](#withmodelvalidations)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Capabilities

- The API will prevent users from performing disallowed actions, adding capabilities is purely to improve UX
- Always test the capability works as expected (never assume the API path 🙂), the extra string interpolation can lead to sneaky typos and incorrect returns from the getters

```js
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default class FooModel extends Model {
  @attr backend;
  @attr('string') fooId;

  // use string interpolation for dynamic parts of API path
  // the first arg is the apiPath, and the rest are the model attribute paths for those values
  @lazyCapabilities(apiPath`${'backend'}/foo/${'fooId'}`, 'backend', 'fooId') fooPath;

  // explicitly check for false because default behavior is showing the thing (i.e. the capability hasn't loaded yet and is undefined)
  get canRead() {
    return this.fooPath.get('canRead') !== false;
  }
  get canEdit() {
    return this.fooPath.get('canUpdate') !== false;
  }
}
```

## Decorators

### [@withFormFields()](../app/decorators/model-form-fields.js)

- Sets `allFields`, `formFields` and/or `formFieldGroups` properties on a model class
- `allFields` includes every model attribute (regardless of args passed to decorator)
- `formFields` and `formFieldGroups` only exist if the relevant arg is passed to the decorator

```js
import { withFormFields } from 'vault/decorators/model-form-fields';

const formFieldAttrs = ['attrName', 'anotherAttr'];
const formGroupObjects = [
  // In form-field-groups.hbs form toggle group names are determined by key names
  // 'default' attribute fields render before any toggle groups
  //  additional attribute fields render inside toggle groups
  { default: ['someAttribute'] },
  { 'Additional options': ['anotherAttr'] },
];

@withFormFields(formFieldAttrs, formGroupObjects)
export default class SomeModel extends Model {
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

### [@withModelValidations()](../app/decorators/model-validations.js)

- Adds `validate()` method on model to check attributes are valid before making an API request
- Option to write a custom validation, or use validation method from the [validators util](../app/utils/validators.js) which is referenced by the `type` key
- Option to add `level: 'warn'` to draw user attention to the input, without preventing form submission
- Component example [here](../lib/pki/addon/components/pki-generate-root.ts)

```js
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  // object key is the model's attribute name
  password: [{ type: 'presence', message: 'Password is required' }],
  keyName: [
    {
      validator(model) {
        return model.keyName === 'default' ? false : true;
      },
      message: `Key name cannot be the reserved value 'default'`,
    },
  ],
};

@withModelValidations(validations)
export default class FooModel extends Model {}

// calling validate() returns an object:
model.validate() = {
  isValid: false,
  state: {
    password: {
      errors: ['Password is required.'],
      warnings: [],
      isValid: false,
    },
    keyName: {
      errors: ["Key name cannot be the reserved value 'default'"],
      warnings: [],
      isValid: true,
    },
  },
  invalidFormMessage: 'There are 2 errors with this form.',
};
```
