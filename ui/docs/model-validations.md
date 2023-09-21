# Model Validations Decorator
  
The model-validations decorator provides a method on a model class which may be used for validating properties based on a provided rule set.

## API

The decorator expects a validations object as the only argument with the following shape:

``` js
const validations = {
  [propertyKeyName]: [
    { type, options, message, level, validator }
  ]
};
```
**propertyKeyName** [string] - each key in the validations object should refer to the property on the class to apply the validation to.
 
**type** [string] - the type of validation to apply. These must be exported from the [validators util](../app/utils/validators.js) for lookup. Type is required if a *validator* function is not provided.

**options** [object] - an optional object for the given validator -- min, max, nullable etc.

**message** [string | function] - string added to the errors array and returned in the state object from the validate method if validation fails. A function may also be provided with the model as the lone argument that returns a string. Since this value is typically displayed to the user it should be a complete sentence with proper punctuation.

**level** [string] *optional* - string that defaults to 'error'. Currently the only other accepted value is 'warn'.

**validator** [function] *optional* - a function that may be used in place of type that is invoked in the validate method. This is useful when specific validations are needed which may be dependent on other class properties.
This function takes the class context (this) as the only argument and returns true or false.

## Usage

Each property defined in the validations object supports multiple validations provided as an array. For example, *presence* and *containsWhiteSpace* can both be added as validations for a string property. 

```js
const validations = {
  name: [
    { type: 'presence', message: 'Name is required.' },
    {
      type: 'containsWhiteSpace',
      message: 'Name cannot contain whitespace.',
    },
  ],
};
```
Decorate the model class and pass the validations object as the argument

```js
import Model, { attr } from '@ember-data/model';
import withModelValidations from 'vault/decorators/model-validations';

const validations = {
  name: [
    { type: 'presence', message: 'Name is required.' },
  ],
};

@withModelValidations(validations)
class SomeModel extends Model {
  @attr name;
}
```

Validations must be invoked using the validate method which is added directly to the decorated class.

```js
const model = await this.store.findRecord('some-model', id);
const { isValid, state, invalidFormMessage } = model.validate();

if (isValid) {
  await model.save();
} else {
  this.formError = invalidFormMessage;
  this.errors = state;
}
```
**isValid** [boolean] - the validity of the full class. If no properties provided in the validations object are invalid this will be true.

**state** [object] - the error state of the properties defined in the validations object. This object is keyed by the property names from the validations object and each property contains an *isValid* and *errors* value. The *errors* array will be populated with messages defined in the validations object when validations fail. Since a property can have multiple validations, errors is always returned as an array.

**invalidFormMessage** [string] - message describing the number of errors currently present on the model class.
 
```js
const { state } = model.validate();
const { isValid, errors } = state[propertyKeyName];
if (!isValid) {
  this.flashMessages.danger(errors.join('. '));
}
```

## Examples
### Basic

```js
const validations = {
  foo: [
    { type: 'presence', message: 'foo is a required field.' }
  ],
};

@withModelValidations(validations)
class SomeModel extends Model { foo = null; }

const model = new SomeModel();
const { isValid, state } = model.validate();

console.log(isValid); // false
console.log(state.foo.isValid); // false
console.log(state.foo.errors); // ['foo is a required field']
```
### Custom validator

```js
const validations = {
  foo: [{
    validator: (model) => model.bar.includes('test') ? model.foo : false,
    message: 'foo is required if bar includes test.'
  }],
};

@withModelValidations(validations)
class SomeModel extends Model {
  foo = false;
  bar = ['foo', 'baz'];
}

const model = new SomeModel();
const { isValid, state } = model.validate();

console.log(isValid); // false
console.log(state.foo.isValid); // false
console.log(state.foo.errors); // ['foo is required if bar includes test.']

model.foo = true;
model.bar.push('test');

console.log(isValid); // true
console.log(state.foo.isValid); // true
console.log(state.foo.errors); // []
```

### Adding class in template based on validation state

All form validation errors must have a red border around them. Add this by adding a conditional class *has-error-border* to the element.

```js
@action
async save() {
  const { isValid, state } = this.model.validate();

  if (isValid) {
    await this.model.save();
  } else {
    this.isNameInvalid = !state.name.isValid;
  }
}
```

```hbs
 <input class="input field {{if this.isNameInvalid 'has-error-border'}}" />
```