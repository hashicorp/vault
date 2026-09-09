/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import V2Form from 'vault/forms/v2/v2-form';

module('Unit | forms/v2 | v2-form', function (hooks) {
  hooks.beforeEach(function () {
    this.mockApi = {
      sys: {
        mountsEnableSecretsEngineRaw: async () => ({ success: true }),
      },
    };

    this.testConfig = {
      name: 'test-form',
      path: '/test/path',
      title: 'Test Form',
      payload: {
        username: '',
        email: '',
        profile: {
          firstName: '',
          lastName: '',
        },
        enabled: false,
      },
      submit: async (api, payload) => {
        return await api.sys.mountsEnableSecretsEngineRaw(payload);
      },
      sections: [
        {
          name: 'section1',
          title: 'User Information',
          fields: [
            {
              name: 'username',
              type: 'TextInput',
              label: 'Username',
              helperText: 'Enter your username',
              validations: [
                { type: 'required', message: 'Username is required' },
                {
                  type: 'minLength',
                  options: { minLength: 3 },
                  message: 'Username must be at least 3 characters',
                },
              ],
            },
            {
              name: 'email',
              type: 'TextInput',
              label: 'Email',
              validations: [
                { type: 'required', message: 'Email is required' },
                { type: 'email', message: 'Must be a valid email' },
              ],
            },
            {
              name: 'profile.firstName',
              type: 'TextInput',
              label: 'First Name',
            },
            {
              name: 'profile.lastName',
              type: 'TextInput',
              label: 'Last Name',
            },
          ],
        },
        {
          name: 'section2',
          title: 'Settings',
          fields: [
            {
              name: 'enabled',
              type: 'Toggle',
              label: 'Enable feature',
            },
          ],
        },
      ],
    };
  });

  test('it initializes with config object', function (assert) {
    const form = new V2Form(this.testConfig);

    assert.ok(form, 'form instance created');
    assert.strictEqual(form.config.name, 'test-form', 'config name is set');
    assert.deepEqual(form.payload, this.testConfig.payload, 'payload is initialized');
    assert.ok(form.isValid, 'form is initially valid');
    assert.strictEqual(form.validationErrors.size, 0, 'no validation errors initially');
  });

  test('it creates a deep copy of the payload', function (assert) {
    const form = new V2Form(this.testConfig);

    // Modify the form payload
    form.payload.username = 'test-user';

    // Original config payload should be unchanged
    assert.strictEqual(this.testConfig.payload.username, '', 'original payload unchanged');
    assert.strictEqual(form.payload.username, 'test-user', 'form payload updated');
  });

  test('it updates simple field values with set()', function (assert) {
    const form = new V2Form(this.testConfig);

    form.set('username', 'john_doe');

    assert.strictEqual(form.payload.username, 'john_doe', 'username updated');
    assert.strictEqual(form.payload.email, '', 'other fields unchanged');
  });

  test('it updates nested field values with dotted path notation', function (assert) {
    const form = new V2Form(this.testConfig);

    form.set('profile.firstName', 'John');
    form.set('profile.lastName', 'Doe');

    assert.strictEqual(form.payload.profile.firstName, 'John', 'nested firstName updated');
    assert.strictEqual(form.payload.profile.lastName, 'Doe', 'nested lastName updated');
  });

  test('it updates nested field values when intermediate objects exist', function (assert) {
    const config = {
      ...this.testConfig,
      payload: {
        username: '',
        user: {
          profile: {
            name: '',
          },
        },
      },
      sections: [
        {
          name: 'section1',
          fields: [
            {
              name: 'user.profile.name',
              type: 'TextInput',
              label: 'Name',
            },
          ],
        },
      ],
    };

    const form = new V2Form(config);
    form.set('user.profile.name', 'John');

    assert.strictEqual(form.payload.user.profile.name, 'John', 'nested value set correctly');
  });

  test('it validates field on set()', function (assert) {
    const form = new V2Form(this.testConfig);

    // Set invalid value (too short)
    form.set('username', 'ab');

    assert.notOk(form.isValid, 'form is invalid');
    assert.strictEqual(form.validationErrors.size, 1, 'has one validation error');
    assert.ok(form.getErrors('username').length > 0, 'username has validation errors');
  });

  test('it clears validation errors when field becomes valid', function (assert) {
    const form = new V2Form(this.testConfig);

    // Set invalid value
    form.set('username', 'ab');
    assert.notOk(form.isValid, 'form is invalid');

    // Set valid value
    form.set('username', 'john_doe');
    assert.ok(form.isValid, 'form is valid');
    assert.strictEqual(form.validationErrors.size, 0, 'no validation errors');
  });

  test('getErrors() returns errors for a field', function (assert) {
    const form = new V2Form(this.testConfig);

    form.set('username', ''); // Required field

    const errors = form.getErrors('username');
    assert.ok(errors.length > 0, 'has errors');
    assert.ok(errors[0].includes('required'), 'error message mentions required');
  });

  test('getErrors() returns empty array for valid field', function (assert) {
    const form = new V2Form(this.testConfig);

    form.set('username', 'valid_username');

    const errors = form.getErrors('username');
    assert.strictEqual(errors.length, 0, 'no errors for valid field');
  });

  test('validateForm() validates all fields', function (assert) {
    const form = new V2Form(this.testConfig);

    // Leave required fields empty
    const result = form.validateForm();

    assert.notOk(result.isValid, 'form is invalid');
    assert.notOk(form.isValid, 'isValid property is false');
    assert.ok(form.validationErrors.size > 0, 'has validation errors');
    assert.ok(form.getErrors('username').length > 0, 'username has errors');
    assert.ok(form.getErrors('email').length > 0, 'email has errors');
  });

  test('validateForm() returns valid when all fields are valid', function (assert) {
    const form = new V2Form(this.testConfig);

    form.set('username', 'john_doe');
    form.set('email', 'john@example.com');

    const result = form.validateForm();

    assert.ok(result.isValid, 'form is valid');
    assert.ok(form.isValid, 'isValid property is true');
    assert.strictEqual(form.validationErrors.size, 0, 'no validation errors');
  });

  test('submit() validates form before submission', async function (assert) {
    const form = new V2Form(this.testConfig);

    // Try to submit with invalid data
    try {
      await form.submit(this.mockApi);
      assert.ok(false, 'should have thrown error');
    } catch (error) {
      assert.ok(error.message.includes('validation failed'), 'throws validation error');
    }
  });

  test('submit() calls config submit handler with valid data', async function (assert) {
    const form = new V2Form(this.testConfig);

    form.set('username', 'john_doe');
    form.set('email', 'john@example.com');

    const response = await form.submit(this.mockApi);

    assert.ok(response.success, 'submit returns response');
  });

  test('it auto-injects required validation for fields with isRequired: true', function (assert) {
    const config = {
      ...this.testConfig,
      sections: [
        {
          name: 'section1',
          fields: [
            {
              name: 'username',
              type: 'TextInput',
              label: 'Username',
              isRequired: true,
            },
          ],
        },
      ],
    };

    const form = new V2Form(config);

    // Check that required validation was injected
    const usernameField = form.config.sections[0].fields[0];
    assert.ok(usernameField.validations, 'validations array exists');
    assert.ok(
      usernameField.validations.some((v) => v.type === 'required'),
      'required validation was injected'
    );
  });

  test('it does not duplicate required validation if already present', function (assert) {
    const config = {
      ...this.testConfig,
      sections: [
        {
          name: 'section1',
          fields: [
            {
              name: 'username',
              type: 'TextInput',
              label: 'Username',
              isRequired: true,
              validations: [{ type: 'required', message: 'Custom required message' }],
            },
          ],
        },
      ],
    };

    const form = new V2Form(config);

    const usernameField = form.config.sections[0].fields[0];
    const requiredValidations = usernameField.validations.filter((v) => v.type === 'required');
    assert.strictEqual(requiredValidations.length, 1, 'only one required validation exists');
    assert.strictEqual(
      requiredValidations[0].message,
      'Custom required message',
      'keeps original required validation'
    );
  });

  test('it handles conditional field visibility', function (assert) {
    const config = {
      ...this.testConfig,
      sections: [
        {
          name: 'section1',
          fields: [
            {
              name: 'showAdvanced',
              type: 'Toggle',
              label: 'Show Advanced',
            },
            {
              name: 'advancedOption',
              type: 'TextInput',
              label: 'Advanced Option',
              isVisible: (payload) => payload.showAdvanced === true,
            },
          ],
        },
      ],
    };

    const form = new V2Form(config);

    // Initially, advanced option should not be visible
    form.validateForm();
    assert.strictEqual(form.validationErrors.size, 0, 'hidden field not validated');

    // Enable advanced options
    form.set('showAdvanced', true);

    // Now if we add validation to the advanced field, it should be validated
    const advancedField = form.config.sections[0].fields[1];
    advancedField.validations = [{ type: 'required', message: 'Required' }];

    form.validateForm();
    assert.ok(form.validationErrors.size > 0, 'visible field is validated');
  });

  test('it prunes validation errors for hidden fields', function (assert) {
    const config = {
      ...this.testConfig,
      sections: [
        {
          name: 'section1',
          fields: [
            {
              name: 'showAdvanced',
              type: 'Toggle',
              label: 'Show Advanced',
            },
            {
              name: 'advancedOption',
              type: 'TextInput',
              label: 'Advanced Option',
              isVisible: (payload) => payload.showAdvanced === true,
              validations: [{ type: 'required', message: 'Required' }],
            },
          ],
        },
      ],
    };

    const form = new V2Form(config);

    // Enable advanced and trigger validation error
    form.set('showAdvanced', true);
    form.set('advancedOption', ''); // Invalid
    assert.ok(form.getErrors('advancedOption').length > 0, 'has validation error');

    // Disable advanced - error should be pruned
    form.set('showAdvanced', false);
    assert.strictEqual(form.getErrors('advancedOption').length, 0, 'error pruned for hidden field');
  });
});
