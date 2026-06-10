/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupRenderingTest } from 'vault/tests/helpers';

module('Integration | Component | form/v2/field', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.onChange = sinon.spy();
  });

  test('it renders a text input field', async function (assert) {
    this.field = {
      name: 'username',
      label: 'Username',
      type: 'TextInput',
      helperText: 'Enter your username',
    };
    this.value = 'test-user';

    await render(hbs`
      <Form::V2::Field
        @field={{this.field}}
        @value={{this.value}}
        @onChange={{this.onChange}}
      />
    `);

    assert.dom('label').includesText('Username', 'renders label');
    assert.dom('input[type="text"]').exists('renders text input');
    assert.dom('.hds-form-helper-text').includesText('Enter your username', 'renders helper text');
  });

  test('it renders a TextInput email variant', async function (assert) {
    this.field = {
      name: 'email',
      label: 'Email',
      type: 'TextInput',
      inputType: 'email',
    };
    this.value = 'person@example.com';

    await render(hbs`
      <Form::V2::Field
        @field={{this.field}}
        @value={{this.value}}
        @onChange={{this.onChange}}
      />
    `);

    assert.dom('input[type="email"]').exists('renders email input type');
  });

  test('it renders a TextInput password variant', async function (assert) {
    this.field = {
      name: 'password',
      label: 'Password',
      type: 'TextInput',
      inputType: 'password',
    };

    await render(hbs`
      <Form::V2::Field
        @field={{this.field}}
        @onChange={{this.onChange}}
      />
    `);

    assert.dom('input[type="password"]').exists('renders password input type');
  });

  test('it renders validation errors', async function (assert) {
    this.field = {
      name: 'password',
      label: 'Password',
      type: 'TextInput',
    };
    this.errors = ['Password is required'];

    await render(hbs`
      <Form::V2::Field
        @field={{this.field}}
        @errors={{this.errors}}
        @onChange={{this.onChange}}
      />
    `);

    assert.dom('.hds-form-error__message').includesText('Password is required', 'renders error message');
  });

  test('it renders a textarea field', async function (assert) {
    this.field = {
      name: 'description',
      label: 'Description',
      type: 'TextArea',
      helperText: 'Enter a detailed description',
    };
    this.value = 'Multi-line\ntext content';

    await render(hbs`
      <Form::V2::Field
        @field={{this.field}}
        @value={{this.value}}
        @onChange={{this.onChange}}
      />
    `);

    assert.dom('label').includesText('Description', 'renders label');
    assert.dom('textarea').exists('renders textarea');
    assert.dom('.hds-form-helper-text').includesText('Enter a detailed description', 'renders helper text');
  });

  test('it renders a toggle field', async function (assert) {
    this.field = {
      name: 'enabled',
      label: 'Enable feature',
      type: 'Toggle',
      helperText: 'Turn this on to enable the feature',
    };
    this.value = true;

    await render(hbs`
      <Form::V2::Field
        @field={{this.field}}
        @value={{this.value}}
        @onChange={{this.onChange}}
      />
    `);

    assert.dom('label').includesText('Enable feature', 'renders label');
    assert.dom('input[type="checkbox"]').exists('renders toggle');
    assert
      .dom('.hds-form-helper-text')
      .includesText('Turn this on to enable the feature', 'renders helper text');
  });

  test('it renders a checkbox field', async function (assert) {
    this.field = {
      name: 'agree',
      label: 'I agree to the terms',
      type: 'Checkbox',
    };
    this.value = false;

    await render(hbs`
      <Form::V2::Field
        @field={{this.field}}
        @value={{this.value}}
        @onChange={{this.onChange}}
      />
    `);

    assert.dom('label').includesText('I agree to the terms', 'renders label');
    assert.dom('input[type="checkbox"]').exists('renders checkbox');
  });

  test('it renders a radio group', async function (assert) {
    this.field = {
      name: 'size',
      label: 'Select size',
      type: 'Radio',
      options: [
        { label: 'Small', value: 'small' },
        { label: 'Medium', value: 'medium' },
        { label: 'Large', value: 'large' },
      ],
    };
    this.value = 'medium';

    await render(hbs`
      <Form::V2::Field
        @field={{this.field}}
        @value={{this.value}}
        @onChange={{this.onChange}}
      />
    `);

    assert.dom('legend').includesText('Select size', 'renders legend');
    assert.dom('input[type="radio"]').exists({ count: 3 }, 'renders radio options');
  });

  test('it renders radio cards', async function (assert) {
    this.field = {
      name: 'plan',
      label: 'Choose a plan',
      type: 'RadioCard',
      options: [
        { label: 'Basic', value: 'basic', description: 'For small teams' },
        { label: 'Pro', value: 'pro', description: 'For growing teams' },
        { label: 'Enterprise', value: 'enterprise', description: 'For large organizations' },
      ],
    };
    this.value = 'pro';

    await render(hbs`
      <Form::V2::Field
        @field={{this.field}}
        @value={{this.value}}
        @onChange={{this.onChange}}
      />
    `);

    assert.dom('legend').includesText('Choose a plan', 'renders legend');
    assert.dom('input[type="radio"]').exists({ count: 3 }, 'renders radio card options');
  });

  test('it renders a masked input field', async function (assert) {
    this.field = {
      name: 'password',
      label: 'Password',
      type: 'MaskedInput',
      helperText: 'Enter a secure password',
    };

    await render(hbs`
      <Form::V2::Field
        @field={{this.field}}
        @onChange={{this.onChange}}
      />
    `);

    assert.dom('label').includesText('Password', 'renders label');
    assert.dom('input').exists('renders masked input');
    assert.dom('.hds-form-helper-text').includesText('Enter a secure password', 'renders helper text');
  });

  test('it renders a select field', async function (assert) {
    this.field = {
      name: 'country',
      label: 'Select country',
      type: 'Select',
      helperText: 'Choose your country',
      options: [
        { label: 'United States', value: 'us' },
        { label: 'Canada', value: 'ca' },
        { label: 'United Kingdom', value: 'uk' },
      ],
    };
    this.value = 'ca';

    await render(hbs`
      <Form::V2::Field
        @field={{this.field}}
        @value={{this.value}}
        @onChange={{this.onChange}}
      />
    `);

    assert.dom('label').includesText('Select country', 'renders label');
    assert.dom('select').exists('renders select element');
    assert.dom('select option').exists({ count: 3 }, 'renders 3 options');
    assert.dom('select').hasValue('ca', 'select has correct value');
    assert.dom('.hds-form-helper-text').includesText('Choose your country', 'renders helper text');
  });

  test('it renders a select field with custom placeholder', async function (assert) {
    this.field = {
      name: 'region',
      label: 'Select region',
      type: 'Select',
      placeholder: 'Pick a region',
      options: [
        { label: 'North America', value: 'na' },
        { label: 'Europe', value: 'eu' },
      ],
    };

    await render(hbs`
      <Form::V2::Field
        @field={{this.field}}
        @onChange={{this.onChange}}
      />
    `);

    assert.dom('select option:first-child').hasText('Pick a region', 'renders custom placeholder');
    assert.dom('select option:first-child').hasValue('', 'placeholder has empty value');
  });

  test('it renders unsupported field types as a text input fallback with console warning', async function (assert) {
    assert.expect(6);

    this.field = {
      name: 'unsupported',
      label: 'Unsupported Field',
      type: 'FutureFieldType',
    };

    // Capture console.warn calls
    const originalWarn = console.warn;
    let warnCalled = false;
    let warnMessage = '';
    console.warn = (message) => {
      warnCalled = true;
      warnMessage = message;
    };

    await render(hbs`
      <Form::V2::Field
        @field={{this.field}}
        @onChange={{this.onChange}}
      />
    `);

    // Restore console.warn
    console.warn = originalWarn;

    // Assert text input is rendered as fallback
    assert.dom('.hds-form-text-input').exists('renders text input as fallback');
    assert.dom('label').hasText('Unsupported Field', 'displays correct label');

    // Assert console warning was logged
    assert.true(warnCalled, 'console.warn was called');
    assert.true(
      warnMessage.includes('Unsupported field type "FutureFieldType"'),
      'warning includes field type'
    );
    assert.true(warnMessage.includes('Unsupported Field'), 'warning includes field label');
    assert.true(warnMessage.includes('Falling back to text input'), 'warning includes fallback message');
  });
});
