# Background

The `Form` class was created as a replacement for form related functionality
that previously lived in Ember Data models. Given that the `FormField` component
was designed around the metadata that was defined on model attributes, it was
imperative to preserve this pattern while moving the functionality to a dependency-free
native javascript solution.

# Usage

The `Form` class is intended to be extended by a class that represents a particular form
in the application.

```ts
export default class MyForm extends Form {
  declare data: MyFormData;

  // define form fields
  name = new FormField('name', 'string');
  secret = new FormField('secret', 'string', {
    editType: 'kv',
    keyPlaceholder: 'Secret key',
    valuePlaceholder: 'Secret value',
    label: 'Secret (kv pairs)',
    isSingleRow: true,
    allowWhiteSpace: true,
  });

  // define validations
  validations: Validations = {
    name: [{ type: 'presence', message: 'Name is required.' }],
  };

  // if serialization is needed override toJSON method
  toJSON() {
    const trimmedName = this.data.name.trim();
    return super.toJSON({ ...this.data, name: trimmedName });
  }
}
```

Form data is set to the data object on the class and can be initialized
with defaults or server data when editing by passing an object into the constructor.

```ts
// create route
model() {
  return new MyForm({ name: 'Default name' });
}

// edit route
async model() {
  const data = await this.api.fetchSomeData();
  return new MyForm(data);
}
```

The route model (`MyForm` instance) can be passed into the form component in
the same manner as an Ember Data model and the `formFields` can be looped
to render `FormField` components.

```hbs
{{#each @form.formFields as |field|}}
  <FormField @attr={{field}} @model={{@form}} @modelValidations={{this.validations}} />
{{/each}}
```

To validate the form and access the data use the toJSON method.

```ts
// save method of form component
async save() {
  try {
    const { isValid, state, invalidFormMessage, data } = this.args.form.toJSON();
    this.validations = isValid ? null : state;
    this.invalidFormMessage = invalidFormMessage;

    if (isValid) {
      await this.api.saveTheForm(data);
      this.flashMessages.success('It worked');
      this.router.transitionTo('another.route');
    }
  } catch(error) {
    const { message } = await this.api.parseError(error);
    this.errorMessage = message;
  }
}
```
