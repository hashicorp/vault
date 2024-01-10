# Forms

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Forms](#forms)
  - [How to build a form](#how-to-build-a-form)
    - [Step 1 - Define the model](#step-1---define-the-model)
    - [Step 2 - define the sections and fields](#step-2---define-the-sections-and-fields)
    - [Step 3 - Create your form component](#step-3---create-your-form-component)
      - [Custom errors](#custom-errors)
    - [Step 4 - Render the field inputs](#step-4---render-the-field-inputs)
      - [Without FormField](#without-formfield)
      - [With expanded field groups](#with-expanded-field-groups)
    - [Step 5 - Add validation](#step-5---add-validation)
  - [General Patterns](#general-patterns)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## How to build a form

So you want to build a form -- you're in good company! Here's a general guideline on how to build the forms.

### Step 1 - Define the model

Ember-Data model or a native JS class? While both should work, most of the Vault UI forms are based on ember-data models, so we're going to continue with Ember-Data as our example for the rest of this document.

```js
export default class UserModel extends Model {
  @attr('string') firstName;
  @attr('string') lastName;
  @attr('string', {
    possibleValues: ['R&D', 'Sales', 'Marketing']
  }) department;
  @attr('string', {
    editType: 'searchSelect'
  }) manager;
  @attr('boolean') onboarded;
  @attr('boolean', {
    helpText: 'Check this if you are an admin'
    editType: 'toggle'
  }) admin;
  @attr('string', {
    label: 'Reason user is admin'
  }) adminReason;
}
```

In this model for a user, there are a few fields that are relevant to every user, and some fields whose values depend on other values on the model. Some acceptance criteria for this form: 

- Depending on the value of `department`, the `search-select` will query and show different manager values
- Both `onboarded` and `admin` are boolean values, but since selecting `admin` will show new field options (`adminReason`), the `admin` editType is `toggle` while `onboarded` is a checkbox since there are no UI side effects.
- We can imagine there are other fields that are specific to the department, and will be shown in their own section on the form.

Now that we have a layout of the data working with, let's start to put the form together.

### Step 2 - define the sections and fields

The simplest form will have one set of fields that are shown no matter what. In more complex examples, there are multiple sections whose fields may change depending on model values. That is the case in our example above.

The sections might be represented like this, where the title of the section is the key and the array of field names is the value:

```js
{
  default: ['firstName', 'lastName', 'department'],
  "User details": ['onboarded', 'admin'], // + adminReason, if admin = true
  "Department options": [] // Fully dependent on value of `department`
}
```

Now that we have some understanding of which fields go where, let's get templating!

### Step 3 - Create your form component

Each form will still need its own component, with a backing class to define the save and cancel functionality. However, the EasyForm component makes much of the form boilerplate.

Create your form component `ember g component -gc user-form` and fill in with the following to create a bare minimum form. This component assumes we're passing the ember-data model created above as `@model`.

```js
export default class UserForm extends Component {
  @service router;
  @service flashMessages;

  @action async saveModel() {
    await this.args.model.save();
    this.flashMessages.success('User was saved');
    this.router.transitionTo('user.detail', this.args.model.id);
  }
}
```

```hbs
<EasyForm @onSave={{perform this.saveModel}} @onCancel={{transition-to 'overview'}} as |F|>
  <F.NamespaceReminder @noun='user' @mode={{if @model.isNew 'create' 'edit'}} />
</EasyForm>
```

With just this minimum amount of code, we have a form which:

- Shows a Namespace Reminder if the form is rendered within a namespace
- Has working Save and Cancel buttons at the bottom of the form
- Prevents default form event when form is submitted
- Shows a loading icon when the onSave method is running (even if the save method is not a concurrency task!)
- Shows a formatted error on fail (even though we don't handle the error at all in the `saveModel` function!)

#### Custom errors

What if we want to show a custom error message at some point in the save method? As long as an `Error` is thrown, the EasyForm's submit handler will be able to manage the error message.

```js
// user-form.js
@action saveModel() {
  if (!this.args.model.someField !== 'Accepted') {
    throw new Error("Some field is not accepted!")
  }
  try {
    await this.args.model.save();
    this.flashMessages.success('User was saved');
    this.router.transitionTo('user.detail', this.args.model.id);
  } catch (e) {
    if (e.httpStatus === '403') {
      throw new Error('Permission was denied');
    }
    // as a fallback, throw the original error
    throw e;
  }
}
```

### Step 4 - Render the field inputs

Now we're finally going handle the meat and potatoes of the form: fields. There are a few ways to go about this but essentially we want to minimize the amount of developer decisions and keep most of the form options in either the model options (if ember-data) or form component (if native JS class).

We already added a bunch of of options into the model attributes in the user Model class. In order to access these attribute options in the component, we need to add a decorator to the Model that sets the expanded attributes to a value on the model called `allByKey`. We're also going to define a few getters that return the fields for each section.

```js
@withExpandedAttributes();
export default class UserModel extends Model {
  get mainFields() {
    return ['firstName', 'lastName', 'department'];
  }
  get userFields() {
    const userFields = ['onboarded', 'admin'];
    if (this.admin) {
      userFields.push('adminReason')
    }
    return userFields;
  }
  get departmentFields() {
    if (!this.department) return [];
    // ... calculate which fields to return based on this.department
  }
```

Now in the template we can add the sections. The first field is the default and has no title, but it does have a divider (which is rendered on the bottom). The second section has a title but no divider. The third section is togglable, which requires a title.

```hbs
<EasyForm @onSave={{perform this.saveModel}} @onCancel={{transition-to 'overview'}} as |F|>
  <F.NamespaceReminder @noun='user' @mode={{if @model.isNew 'create' 'edit'}} />
  <F.Section @hasDivider={{true}}>
    {{#each (F.expand-fields @model.mainFields @model.allByKey) as |attr|}}
      <FormField @model={{@model}} @attr={{attr}} />
    {{/each}}
  </F.Section>
  <F.Section @title='User settings'>
    {{#each (F.expand-fields @model.userFields @model.allByKey) as |attr|}}
      <FormField @model={{@model}} @attr={{attr}} />
    {{/each}}
  </F.Section>
  <F.Section @title='Department options' @toggles={{true}}>
    {{#each (F.expand-fields @model.departmentFields @model.allByKey) as |attr|}}
      <FormField @model={{@model}} @attr={{attr}} />
    {{else}}
      <EmptyState @title='Choose a department to see the fields' />
    {{/each}}
  </F.Section>
</EasyForm>
```

#### Without FormField

The previous example uses the `FormField` component, which expects to be passed an Ember-Data model and an expanded attribute. It then maps to the correct field type and **automatically updates the model**.

If you choose not to use `FormField`, you will need to handle model updates manually and wrap the component in `F.Field` which will manage spacing. Let's see a simplified, manual example:

```hbs
<EasyForm @onSave={{perform this.saveModel}} @onCancel={{transition-to 'overview'}} as |F|>
  <F.NamespaceReminder @noun='user' @mode={{if @model.isNew 'create' 'edit'}} />
  <F.Section>
    <F.Field>
      <FieldInput::TextInput
        @name='firstName'
        @label='First Name'
        @value={{this.userState.firstName}}
        @onChange={{this.handleChange}}
      />
    </F.Field>
    <F.Field>
      <FieldInput::TextInput
        @name='lastName'
        @label='Last Name'
        @value={{this.userState.lastName}}
        @onChange={{this.handleChange}}
      />
    </F.Field>
  </F.Section>
</EasyForm>
```

```js
export default class UserForm extends Component {
  // ...
  @action handleChange(key, value) {
    this.userState[key] = value;
  }
}
```

All of the `FieldInput` components have the same `onChange` callback signature, but the values may have different shapes depending on the type of input. It's as simple as that!

#### With expanded field groups

You may be thinking to yourself, "Why can't my sections be codified too?!" They can! First we need to define the group shapes as another getter on the model. In this example, we're just using the getters we already defined for each section.

```js
export default class UserModel extends Model {
  // ...
  get groups() {
    return [
      { fields: this.mainFields, hasDivider: true },
      { title: 'User settings', fields: this.userFields },
      { title: 'Department options', fields: this.departmentFields, toggles: true },
    ];
  }
}
```

Then in the template, we can iterate over each section as such:

```hbs
<EasyForm @onSave={{perform this.saveModel}} @onCancel={{transition-to 'overview'}} as |F|>
  <F.NamespaceReminder @noun='user' @mode={{if @model.isNew 'create' 'edit'}} />
  {{#each (F.expand-groups @model.groups @model.allByKey) as |group|}}
    <F.Section @hasDivider={{group.hasDivider}} @title={{group.title}} @toggles={{group.toggles}}>
      {{#each group.fields as |attr|}}
        <FormField @model={{@model}} @attr={{attr}} />
      {{/each}}
    </F.Section>
  {{/each}}
</EasyForm>
```

Viola! The downside to this is we no longer have the empty state when departments are empty, so this technique will work in some cases but there will be other times when it is better to have more granular control over the sections.

### Step 5 - Add validation

We add validations on the model level, with a decorator called `withModelValidations`. For presence validation, we have the option of either setting `@isRequired` on the `FieldInput` which will also show a `required` badge on the input, or using a validator. We're going to use the validator in this example, because it's one of the easiest validators to set up. For more on validators, check out [our docs on model-validations](./model-validations.md).

```js
const validations = {
  firstName: [
    { type: 'presence', message: `First name is required` },
  ],
  lastName: [
    { type: 'presence', message: `Last name is required` },
  ]
}
@withModelValidations(validations);
@withExpandedAttributes();
export default class UserModel extends Model {
```

Then we'll want to add the validation method to EasyForm. Note that while it works seamlessly with the `withModelValidations` decorator, it needs only to return `isValid` and `state` in its response to work on the form.

```hbs
<EasyForm
  @onSave={{perform this.saveModel}}
  @onCancel={{transition-to 'overview'}}
  @onValidate={{@model.validate}} as |F|>
  <F.Section>
    <F.Field>
      {{#let (F.is-valid attr.name) as |valid|}}
      <FieldInput::TextInput
        @name='custom'
        @label='My custom field'
        @value={{this.userState.custom}}
        @onChange={{this.handleChange}}
        @fieldErrors={{valid.errors}}
      />
    </F.Field>
    {{#each group.fields as |attr|}}
      <FormField @model={{@model}} @attr={{attr}} @validations={{F.is-valid attr.name}} />
    {{/each}}
  </F.Section>

```

## General Patterns

- Render `FlashMessage` on success
- Handling errors/validation messages:
  - Render API errors using a `<MessageError>` or `Hds::Alert` at the top of forms
  - Display validation error messages `onsubmit` (not `onchange` for inputs)
  - Render an `<AlertInline>` [beside](../lib/pki/addon/components/pki-role-generate.hbs) form buttons, especially if the error banner is hidden from view (long forms). Message options:
    - The `invalidFormMessage` from a model's `validate()` method that includes an error count
    - Generic message for API errors or forms without model validations: 'There was an error submitting this form.'
  - Add `has-error-border` class to invalid inputs
