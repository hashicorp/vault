# Models

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Models](#models)
  - [Intro](#intro)
  - [Model patterns overview](#model-patterns-overview)
  - [Patterns](#patterns)
    - [Attributes \& field groups](#attributes--field-groups)
      - [Attributes \& field groups example](#attributes--field-groups-example)
    - [Validations](#validations)
      - [@withModelValidations()](#withmodelvalidations)
    - [Capabilities](#capabilities)
      - [Examples](#examples)
    - [Models hydrated by OpenAPI](#models-hydrated-by-openapi)
  - [Using Decorators](#using-decorators)
    - [@withFormFields()](#withformfields)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Intro

We use models primarily as the backing data layer for our forms and for our list/show views. As Ember-Data has matured, our patterns of usage have become outdated. This document serves to outline our current best-practices, since examples within the codebase are often out of date and do not always reflect our best practices or ambitions.

## Model patterns overview

Models can be thought of as the shape of data that an instance of that Model -- a Record -- will have. Models should be as "thin" as possible, holding only data directly relevant to the Record itself. For example, if we have a Model `user` with attributes `firstName` and `lastName`, it _is_ appropriate to have a getter on the Model called `fullName`, because its attributes can be calculated directly from the record's values, and is relevant to the Record itself. However it is _not_ appropriate to store data like which fields are shown on the edit form, because that has no bearing on the Record itself. Field values are a display concern, not related to the values of the record.

Other patterns and where they belong in relation to the Model:

- **Attribute metadata** - this is referring to information defined on a Model's attributes, such as label, edit type, and other information relevant to both forms and the given attribute. We use these heavily in the `FormField` component to show the correct label, help text, and input type. Conceptually, this does not belong on a Model (because the information is not directly related to the data in a Record) but, since we leverage OpenAPI heavily to populate both attributes and their metadata, we are going to keep attribute metadata defined on the attribute in the Model. **TL;DR: Lives on Model**

- **Form and show fields** - the grouping and order of fields that should display on both show routes and create/edit forms, while conceptually related to the Model, is not related to an individual record. Therefore, this information should not be defined on the Model (which has been our previous pattern). To support migration, we have a few helpful decorators and patterns. **TL;DR: Lives in component or model-helper util files**

- **Validations** - While an argument can go either way about this one, we are going to continue defining these on the Model using our handy [withModelValidation decorator](#withmodelvalidations). The state of validation is directly correlated to a given Record which is a strong reason to keep it on the Model. **TL;DR: Lives on Model**

- **[Capabilities](#capabilities)** - Capabilities are calculated by fetching permissions by path -- often multiple paths, based on the same information we need to fetch the Record data (eg. backend, ID). When using `lazyCapabilities` on the model we kick off one API request for each of the paths we need, while using the capabilities service `fetchMultiplePaths` method we can make one request with all the required paths included. Our best practice is to fetch capabilities outside of the Model (perhaps as part of a route model, or on a user action such as dropdown click open). A downside to this approach is that the API may get re-requested per page we check the capability (eg. list view dropdown and then detail view) -- but we can optimize this in the future by first checking the store for `capabilities` of matching path/ID before sending the API request. **TL;DR: Lives in route or component where they are used**

## Patterns

### Attributes & field groups

We use attributes defined on the Model to determine input concerns (label, input type, help text) and field groups to determine the order of the attribute data on the form and detail pages, and are defined in the component they are used in or in a `utils/model-helpers/*` file.

#### Attributes & field groups example

In this example, we have a Model `simple-timer` with a few attributes defined. The `withExpandedAttributes` helper adds a couple items to the Model it's applied to:

- allByKey - a getter which returns all the attributes as keys of an object, and the value is the metadata of the attribute including anything returned from OpenAPI if the model is included in `OPENAPI_POWERED_MODELS`.
- \_expandGroups - takes an array of group objects and expands the attribute keys into the metadata

In the component where we pass a Record of this Model, we can see how we use it to populate either a flat array of attributes for use in the show view, or to populate groups of fields for rendering on a form.

```js
// models/simple-timer.js
@withExpandedAttributes()
export default class SimpleTimer extends Model {
  @attr('string', {
    editType: 'ttl',
    defaultValue: '3600s',
    label: 'TTL',
    helpText: 'Here is some help text',
  })
  ttl;

  @attr('string') name;
  @attr('boolean') restartable; // enterprise only
}
```

```js
// components/simple-timer-display.ts
export default class SimpleTimerDisplay extends Component<Args> {
  @service declare readonly version: VersionService;

  // these fields are shown flat in the show mode, iterated over
  // and used in InfoTableRow
  get showFields() {
    let fields = ['name', 'ttl'];
    if (this.version.isEnterprise) {
      fields.push('restartable');
    }
    return fields.map((field) => this.args.model.allByKey[field]);
  }

  // these fields are shown grouped in edit mode and is formatted
  // to be used in something like FormFieldGroups
  get fieldGroups() {
    let groups = [{ default: ['name', 'ttl'] }];
    if (this.version.isEnterprise) {
      groups[{ 'Custom options': ['restartable'] }];
    }
    return this.args.model._expandGroups(groups);
  }
}
```

### Validations

Validations on used on forms, to present the user with feedback about their form answers before sending the payload to the API. Our best practices are:

- define the validations using the `withModelValidations` decorator
- trigger the `validate()` method added by the decorator on form submit
- if there are validation errors:
  - show a message at the bottom of the form saying there were errors with the form
  - add inline-alert next to the inputs that have incorrect data
  - exit the form submit function early
  - do not disable the submit button
- if there are no validation errors, continue saving as normal

#### [@withModelValidations()](../app/decorators/model-validations.js)

This decorator:

- Adds `validate()` method on model to check attributes are valid before making an API request
- Provides option to write a custom validation, or use validation method from the [validators util](../app/utils/model-helpers/validators.js) which is referenced by the `type` key
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
export default class FooModel extends Model {
  @attr() password;
  @attr() keyName;
}
```

```js
// form-component.js
export default class FormComponent extends Component {
  @tracked modelValidations = null;
  @tracked invalidFormAlert = '';

  checkFormValidity() {
    interface Validity {
      // only true if all of the state's isValid are also true
      isValid: boolean;
      state: {
        // state is keyed by the attribute names
        [key: string]: {
          errors: string[];
          warnings: string[];
          isValid: boolean;
        }
      }
      invalidFormMessage: string; // eg "There are 2 errors with this form"
    }
    // calling validate() returns Validity
    const { isValid, state, invalidFormMessage } = this.args.model.validate();
    this.modelValidations = state;
    this.invalidFormAlert = invalidFormMessage;
    return isValid;
  }

  @action
  submit() {
    // clear errors
    this.modelValidations = null;
    this.invalidFormAlert = null;

    // check validity
    const continueSave = this.checkFormValidity();
    if (!continueSave) return;

    // continue save ...
  }
}
```

### Capabilities

- The API will prevent users from performing disallowed actions, so adding capabilities is purely to improve UX by hiding actions we know the user cannot take. Because of this, we default to showing items if we cannot determine the capabilities for an endpoint.
- Always test the capability works as expected (never assume the API path 🙂) -- the extra string interpolation can lead to sneaky typos and incorrect returns from the getters
- Capabilities are checked via the `capabilities-self` endpoint, and registered in the store as a [capabilities Model](../app/models/capabilities.js), with the path as the Record's ID.
- The path IDs on the capabilities Records should never include the namespace, but when operating within a namespace the paths in the API request payload must be prepended with the namespace so the API will return the proper capabilities (eg. for `kv/data/foo` in the `admin` namespace instead of root)
- In general we want to check capabilities outside of the Model, but we have a patterns for both ways.

#### Examples

**Single capability check within a component**
In [this example](../app/components/clients/page-header.js), we have an action that some users can take within the page header. Honestly this capability check could just have easily lived in the route's Model (since the PageHeader always renders on the relevant routes), but here it provides a good example of a check happening on component instantiation, using the args passed to the component:

```js
// clients/page-header.js
constructor() {
  super(...arguments);
  this.getExportCapabilities(this.args.namespace);
}

async getExportCapabilities(ns = '') {
  try {
    const url = ns
      ? `${sanitizePath(ns)}/sys/internal/counters/activity/export`
      : 'sys/internal/counters/activity/export';
    const cap = await this.store.findRecord('capabilities', url);
    this.canDownload = cap.canSudo;
  } catch (e) {
    // if we can't read capabilities, default to show
    this.canDownload = true;
  }
}
```

**Multiple capabilities checked at once**
When there are multiple capabilities paths to check, the recommended approach is to use the [capabilities service's](../app/services/capabilities.ts) `fetchMultiplePaths` method. It will pass all the paths in a single API request instead of making a capabilities-self call for each path as the other techniques do. In [this example](../lib/kv/addon/routes/secret.js), we get the capabilities as part of the route's model hook and then return the relevant `can*` values:

```js
async fetchCapabilities(backend, path) {
  const metadataPath = `${backend}/metadata/${path}`;
  const dataPath = `${backend}/data/${path}`;
  const subkeysPath = `${backend}/subkeys/${path}`;
  const perms = await this.capabilities.fetchMultiplePaths([metadataPath, dataPath, subkeysPath]);
  // returns values keyed at the path
  return {
    metadata: perms[metadataPath],
    data: perms[dataPath],
    subkeys: perms[subkeysPath],
  };
}

async model() {
  const backend = this.secretMountPath.currentPath;
  const { name: path } = this.paramsFor('secret');
  const capabilities = await this.fetchCapabilities(backend, path);
  return hash({
    // ...
    canUpdateData: capabilities.data.canUpdate,
    canReadData: capabilities.data.canRead,
    canReadMetadata: capabilities.metadata.canRead,
    canDeleteMetadata: capabilities.metadata.canDelete,
    canUpdateMetadata: capabilities.metadata.canUpdate,
  });
}
```

Lastly, we have an example that is common but a pattern that we want to move away from: using `lazyCapabilities` on a Model. The `lazyCapabilities` macro only fetches the capabilities when the attribute is invoked -- so in the example below, only when `canRead` is rendered on the template will the capablities-self call be kicked off.

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

### Models hydrated by OpenAPI

In a Model which is hydrated by OpenAPI, it can be cumbersome to keep up with all the changes made by the backend. One pattern available to us is the [`combineFieldGroups`](../app/utils/openapi-to-attrs.js) method, which

---

## Using Decorators

### [@withFormFields()](../app/decorators/model-form-fields.js)

- Sets `allFields`, `formFields` and/or `formFieldGroups` properties on a model class
- `allFields` includes every model attribute (regardless of args passed to decorator)
- `formFields` and `formFieldGroups` only exist if the relevant arg is passed to the decorator
- `type` of validator should match the keys in [model-helpers/validators.js](../app/utils/model-helpers/validators.js)

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
