# V2 Form System

The V2 form system is a data-driven approach to building forms on top of the HashiCorp Design System (HDS). A `FormConfig` object describes the form's fields, sections, and submit logic. The `V2Form` class manages runtime state (payload and validation errors). Ember components render the form and handle submission.

Use this system for new forms connected to the Vault API. If a matching OpenAPI operation exists, generate a typed scaffold with `pnpm generate:form-config` rather than writing the config by hand.

---

## Contents

- [Prerequisites](#prerequisites)
- [Architecture overview](#architecture-overview)
- [Generating a form config](#generating-a-form-config)
- [Output](#output)
- [After generating](#after-generating)
- [Overriding a generated config](#overriding-a-generated-config)
- [Writing a config manually](#writing-a-config-manually)
- [Using a form in a template](#using-a-form-in-a-template)
- [Multi-step wizards](#multi-step-wizards)
- [Validation](#validation)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

- `pnpm install` has been run (installs `tsx` and other dev dependencies)
- The `@hashicorp/vault-client-typescript` package is installed (provides `openapi.json` and typed SDK methods)

---

## Architecture overview

```
FormConfig (plain object)
  └── sections[]
        └── fields[]          ← field name, type, label, validations, options…

V2Form (class, @tracked)
  ├── payload                 ← deep clone of FormConfig.payload
  ├── validationErrors        ← Map<fieldName, string[]>
  ├── set(path, value)        ← updates payload + re-validates field
  ├── validateForm()          ← validates all visible fields
  └── submit(api)             ← validateForm() → config.submit(api, payload)

Ember components
  ├── Form::V2               ← submit task, error state, yields Form + submitTask
  ├── Form::V2::Renderer     ← <Hds::Form>, iterates sections/fields
  ├── Form::V2::Section      ← wraps fields in <Form.Section> with optional title
  ├── Form::V2::Field        ← renders the correct HDS input component
  ├── Form::V2::ErrorAlert   ← inline critical alert for submission errors
  ├── Form::V2::Wizard       ← multi-step wizard orchestrator
  └── Form::V2::Apply        ← final "apply changes" step with code snippet options
```

---

## Generating a form config

Pass the camelCase API method name to the generator. It reads field definitions from the bundled OpenAPI spec and writes a fully-typed scaffold.

```
pnpm generate:form-config mountsEnableSecretsEngine
```

**What happens:**

1. Searches `openapi.json` for the `POST` operation whose `operationId` matches the dasherized method name
2. Extracts path parameters and request body properties, skipping deprecated fields
3. Groups fields using the `x-vault-displayAttrs.group` annotation from the spec (falls back to a `default` group)
4. Derives the TypeScript request type from the operation's tag and method name (e.g., `SystemApiMountsEnableSecretsEngineOperationRequest`)
5. Writes a `.ts` file to `app/forms/v2/generated/` and runs Prettier on it

The API class used in `submit` is determined automatically from the OpenAPI tag:

| Tag | API class |
|-----|-----------|
| `system` | `api.sys` |
| `auth` | `api.auth` |
| `identity` | `api.identity` |
| `secrets` | `api.secrets` |

---

## Output

The generator writes to `app/forms/v2/generated/` with a dasherized filename:

```
app/forms/v2/generated/mounts-enable-secrets-engine-config.ts
```

The file contains:

- An `import` for the typed SDK request type
- A `FormConfig<RequestType, unknown>` constant with `name`, `path`, `description`, `submit`, `payload`, and `sections`
- `sections` grouped as the OpenAPI spec describes (`params` for path parameters, `default` for body fields, named groups for any `x-vault-displayAttrs.group` annotations)
- All fields typed as `TextInput` by default — update `type` after reviewing

Generated files are **not registered automatically**. See [After generating](#after-generating).

---

## After generating

Generated files are marked `⚠️ AUTO-GENERATED FILE - DO NOT EDIT`. Leave them as-is — all UI-specific adjustments belong in an override (see [Overriding a generated config](#overriding-a-generated-config)).

1. **Register the config** — add an export to `app/forms/v2/generated/index.ts`:

   ```ts
   import mountsEnableSecretsEngineConfig from './mounts-enable-secrets-engine-config';

   const GENERATED_CONFIGS = {
     mountsEnableSecretsEngine: mountsEnableSecretsEngineConfig,
   };

   export default GENERATED_CONFIGS;
   ```

2. **Create an override** to adjust field types, labels, visibility, or remove irrelevant fields (see [Overriding a generated config](#overriding-a-generated-config))
3. **Use the form** — instantiate `V2Form` with the registry key or config object (see [Using a form in a template](#using-a-form-in-a-template))

---

## Overriding a generated config

Generated configs model the API surface and should not be edited directly. To adjust field order, labels, visibility rules, or add non-API fields, create an override using `configBuilder`.

Create a file in `app/forms/v2/overrides/`:

```ts
// app/forms/v2/overrides/mounts-enable-secrets-engine-config.ts
import generatedConfig from '../generated/mounts-enable-secrets-engine-config';
import { configBuilder } from './override-field';

export default configBuilder(generatedConfig)
  .removeField('default', 'MountsEnableSecretsEngineRequest.seal_wrap')
  .updateField('params', 'path', {
    label: 'Mount path',
    helperText: 'The path to mount to. Example: "aws/east"',
    isRequired: true,
  })
  .addSection({
    name: 'engine_selection',
    title: 'Engine type',
    fields: [
      {
        name: 'MountsEnableSecretsEngineRequest.type',
        type: 'Select',
        label: 'Type',
        options: [
          { label: 'KV', value: 'kv' },
          { label: 'AWS', value: 'aws' },
        ],
      },
    ],
  }, 0) // position 0 inserts before all other sections
  .build();
```

Register the override in `app/forms/v2/overrides/index.ts`:

```ts
import mountsEnableSecretsEngineConfig from './mounts-enable-secrets-engine-config';

const OVERRIDE_CONFIGS = {
  mountsEnableSecretsEngine: mountsEnableSecretsEngineConfig,
};

export default OVERRIDE_CONFIGS;
```

`getFormConfig()` checks overrides before generated configs, so the override takes precedence automatically.

### Available builder methods

| Method | Description |
|--------|-------------|
| `addSection(section, position?)` | Add a new section (appended by default, or at `position`) |
| `removeSection(sectionName)` | Remove a section by name |
| `updateSection(sectionName, updates)` | Update section `title`, `description`, or `isVisible` |
| `addField(sectionName, field)` | Add a field to an existing section |
| `updateField(sectionName, fieldName, overrides)` | Update any field property except `name` |
| `removeField(sectionName, fieldName)` | Remove a field from a section |
| `moveField(fieldName, fromSection, toSection, position?)` | Move a field between sections |
| `reorderFields(sectionName, fieldNames)` | Reorder fields within a section |
| `build()` | Returns the final `FormConfig` |

---

## Writing a config manually

For forms not backed by an OpenAPI operation (or when the generated scaffold is not useful), write the config directly:

```ts
import type ApiService from 'vault/services/api';
import type { FormConfig } from 'vault/forms/v2/form-config';

interface MyPayload {
  name: string;
  ttl: number;
  enabled: boolean;
}

const myFormConfig: FormConfig<MyPayload, unknown> = {
  name: 'myForm',
  path: '/sys/example/{name}',
  title: 'Create example',
  payload: {
    name: '',
    ttl: 0,
    enabled: false,
  },
  submit: async (api: ApiService, payload: MyPayload) => {
    return await api.sys.someMethodRaw(payload);
  },
  onSuccess: (response) => {
    // optional: redirect, show toast, etc.
  },
  sections: [
    {
      name: 'basic',
      title: 'Basic settings',
      fields: [
        {
          name: 'name',
          type: 'TextInput',
          label: 'Name',
          isRequired: true,
        },
        {
          name: 'ttl',
          type: 'TextInput',
          inputType: 'number',
          label: 'TTL',
          helperText: 'Time-to-live in seconds',
        },
        {
          name: 'enabled',
          type: 'Toggle',
          label: 'Enable',
        },
      ],
    },
  ],
};
```

### Field types

| `type` | HDS component rendered |
|--------|----------------------|
| `TextInput` | `Hds::Form::TextInput::Field` |
| `TextArea` | `Hds::Form::Textarea::Field` |
| `Select` | `Hds::Form::Select::Field` |
| `Toggle` | `Hds::Form::Toggle::Field` |
| `Checkbox` | `Hds::Form::Checkbox::Field` |
| `Radio` | `Hds::Form::Radio::Group` |
| `RadioCard` | `Hds::Form::RadioCard::Group` |
| `MaskedInput` | `Hds::Form::MaskedInput::Field` |

Support for new field types is added on an as-needed basis. If a form you're migrating to V2 requires a component not listed above, add support for it in `Form::V2::Field` at that time. Using an unrecognised `type` falls back to a text input and logs a console warning.

`TextInput` accepts an optional `inputType` property for HTML input types (`number`, `email`, `url`, `password`, etc.).

`Select`, `Radio`, and `RadioCard` require an `options` array:

```ts
options: [
  { label: 'Option A', value: 'a' },
  { label: 'Option B', value: 'b', description: 'Only for RadioCard' },
]
```

### Conditional visibility

Fields and sections accept an `isVisible` property:

```ts
// Static — always hidden
{ isVisible: false }

// Dynamic — based on current payload
{ isVisible: (payload) => payload.type === 'advanced' }
```

Hidden fields are automatically excluded from validation and their errors are cleared on payload change.

---

## Using a form in a template

### Route / component setup

```ts
// my-route.ts or my-component.ts
import V2Form from 'vault/forms/v2/v2-form';

// Option 1: registry key (config must be registered in generated/index.ts or overrides/index.ts)
form = new V2Form('mountsEnableSecretsEngine');

// Option 2: direct config object
form = new V2Form(myFormConfig);
```

### Default usage (auto-rendered fields + submit button)

```hbs
<Form::V2 @form={{this.form}} @onSuccess={{this.handleSuccess}} />
```

### Custom submit button

```hbs
<Form::V2 @form={{this.form}} @onSuccess={{this.handleSuccess}} as |Form submitTask|>
  <Form.Section>
    <Hds::ButtonSet>
      <Hds::Button
        @text="Save"
        @color="primary"
        type="submit"
        disabled={{or (not this.form.isValid) submitTask.isRunning}}
        {{on "click" (perform submitTask)}}
      />
      <Hds::Button @text="Cancel" @color="secondary" @route="vault.cluster.index" />
    </Hds::ButtonSet>
  </Form.Section>
</Form::V2>
```

### Hiding auto-rendered fields (custom layout)

Pass `@hideFields={{true}}` to suppress the auto-rendered fields while keeping submission, error handling, and the `Form` context available.

```hbs
<Form::V2 @form={{this.form}} @hideFields={{true}} as |Form submitTask|>
  {{! render fields manually using Form.Section, Form.Field, etc. }}
</Form::V2>
```

---

## Multi-step wizards

Define a `WizardConfig` and pass it to `Form::V2::Wizard`:

```ts
import type { WizardConfig } from 'vault/forms/v2/form-config';
import step1Config from 'vault/forms/v2/generated/step-one-config';
import step2Config from 'vault/forms/v2/generated/step-two-config';

const wizardConfig: WizardConfig = {
  title: 'Enable secrets engine',
  applyChanges: true, // adds a final "Apply changes" step with code snippet options
  steps: [
    {
      name: 'mountConfig',
      title: 'Mount configuration',
      heading: 'Configure the mount',
      formConfig: step1Config,
    },
    {
      name: 'engineConfig',
      title: 'Engine settings',
      formConfig: {
        ...step2Config,
        // Dynamic payload: read the path entered in step 1
        payload: (wizardState) => ({
          ...step2Config.payload,
          mount: wizardState.mountConfig?.payload?.path ?? '',
        }),
      },
    },
  ],
};
```

```hbs
<Form::V2::Wizard
  @config={{this.wizardConfig}}
  @onSuccess={{this.handleComplete}}
  @onCancel={{this.handleCancel}}
/>
```

**Cross-step data sharing:** Step payloads can be functions `(wizardState) => payload`. `wizardState` is a map keyed by step `name`, each containing `{ payload, response, error? }` for completed steps. This lets later steps pre-populate fields from earlier step responses.

**`applyChanges`:** When `true`, the wizard appends a final "Apply changes" step that renders `Form::V2::Apply` — a radio card chooser letting the user apply via Terraform HCL, Vault CLI/API curl commands, or directly through the UI.

---

## Validation

### Built-in validators

Add `validations` to any field:

```ts
{
  name: 'email',
  type: 'TextInput',
  label: 'Email',
  isRequired: true, // shorthand — auto-injects a 'required' rule
  validations: [
    { type: 'email', message: 'Please enter a valid email address' },
    { type: 'maxLength', message: 'Email must be under 255 characters', options: { maxLength: 255 } },
  ],
}
```

| `type` | Validates |
|--------|-----------|
| `required` | Non-empty value (rejects `null`, `undefined`, `''`, `[]`, `{}`) |
| `email` | Email format |
| `url` | URL format (uses `new URL()`) |
| `pattern` | Regex — provide `options.pattern` (string or `RegExp`) and optional `options.flags` |
| `minLength` | Minimum string length — provide `options.minLength` |
| `maxLength` | Maximum string length — provide `options.maxLength` |
| `min` | Minimum numeric value — provide `options.min` |
| `max` | Maximum numeric value — provide `options.max` |

### Custom validators

```ts
{
  name: 'path',
  type: 'TextInput',
  label: 'Path',
  validations: [
    {
      validator: (formData) => !String(formData.path).includes(' '),
      message: 'Path must not contain spaces',
    },
  ],
}
```

Custom validators receive the entire form payload as `formData`, so cross-field validation is possible.

### `isRequired` shorthand

Setting `isRequired: true` on a field automatically prepends a `required` validation rule. It also passes `@isRequired` to the HDS component (renders an asterisk on the label).

### Validation timing

- Fields validate **on change** (after each `form.set()` call)
- All fields validate on **submit** (before calling `config.submit`)
- Hidden fields are excluded from validation automatically

---

## Troubleshooting

**`Form configuration not found for: "…"`** — the config key is not registered. Add it to `app/forms/v2/generated/index.ts` or `app/forms/v2/overrides/index.ts`.

**`Operation "…" not found in openapi.json`** — the method name may be misspelled or the operation may be enterprise-only / not yet in the bundled spec. Check `node_modules/@hashicorp/vault-client-typescript/openapi.json` for the exact `operationId`.

**`Could not determine API class for "…"`** — the OpenAPI operation is missing a recognized tag (`system`, `auth`, `identity`, `secrets`). Write the config manually instead.

**`[Form::V2::Field] Unsupported field type "…"`** — the `type` value in the config does not match a supported `FormElement`. Check the [field types table](#field-types) and update the config.

**`Section "…" not found`** — a `configBuilder` method was called with a section name that does not exist in the base config. Call `builder.getSections()` to inspect available section names.

**Fields not validating** — ensure `isRequired: true` or a `validations` array is present on the field. Fields with no validations are always considered valid.
