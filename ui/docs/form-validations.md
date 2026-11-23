# Form Validations

Forms validate their data with a rule set passed to the [`Form`](../app/forms/form.ts) class. Calling `form.toJSON()` runs the rules, so validating and building the request payload happen in one step.

## API

Validations are an object with the following shape, defined as the `validations` property of a `Form` subclass or passed as the third constructor argument when the rules depend on how the form is created.

```js
const validations = {
  [propertyKeyName]: [{ type, options, message, level, validator }],
};
```

**propertyKeyName** [string] - each key refers to a property of the form's `data`. Dot notation may be used for nested properties.

**type** [string] - the type of validation to apply. These must be exported from the [validators util](../app/utils/forms/validators.js) for lookup. Type is required if a _validator_ function is not provided.

**options** [object] - an optional object for the given validator -- min, max, nullable etc.

**message** [string | function] - string added to the errors array and returned in the state object if validation fails. A function may also be provided with the form data as the lone argument that returns a string. Since this value is typically displayed to the user it should be a complete sentence with proper punctuation.

**level** [string] _optional_ - string that defaults to 'error'. Currently the only other accepted value is 'warn'. Warnings are returned in the state object but do not make the form invalid.

**validator** [function] _optional_ - a function that may be used in place of type. It is called with the form data and the `options`, and returns true or false. This is useful when a rule depends on other properties.

Every field with `editType: 'ttl'` also gets a `maxDuration` rule automatically so durations Vault cannot store are rejected.

## Usage

Each property supports multiple validations provided as an array. For example, _presence_ and _containsWhiteSpace_ can both be added to a string property.

```js
export default class SomeForm extends Form {
  validations = {
    name: [
      { type: 'presence', message: 'Name is required.' },
      { type: 'containsWhiteSpace', message: 'Name cannot contain whitespace.' },
    ],
  };
}
```

`toJSON()` returns the validation result together with the form data:

```js
const { isValid, state, invalidFormMessage, data } = this.form.toJSON();

if (isValid) {
  await this.api.secrets.someCreate(this.form.data.name, data);
} else {
  this.formError = invalidFormMessage;
  this.validationErrors = state;
}
```

**isValid** [boolean] - the validity of the full form. If no properties in the validations object are invalid this will be true.

**state** [object] - the error state of the properties defined in the validations object. This object is keyed by the property names from the validations object and each property contains _isValid_, _errors_ and _warnings_ values. Since a property can have multiple validations, errors is always returned as an array.

**invalidFormMessage** [string] - message describing the number of errors currently present on the form.

**data** [object] - the form data, which subclasses may override `toJSON` to serialize.

```js
const { state } = this.form.toJSON();
const { isValid, errors } = state[propertyKeyName];
if (!isValid) {
  this.flashMessages.danger(errors.join('. '));
}
```

## Examples

### Custom validator

```js
validations = {
  foo: [
    {
      validator: (data) => (data.bar.includes('test') ? !!data.foo : true),
      message: 'foo is required if bar includes test.',
    },
  ],
};
```

### Adding class in template based on validation state

All form validation errors must have a red border around them. Add this by adding a conditional class _has-error-border_ to the element.

```js
@action
async save() {
  const { isValid, state, data } = this.form.toJSON();

  if (isValid) {
    await this.api.secrets.someCreate(data);
  } else {
    this.isNameInvalid = !state.name.isValid;
  }
}
```

```hbs
<input class='input field {{if this.isNameInvalid "has-error-border"}}' />
```
