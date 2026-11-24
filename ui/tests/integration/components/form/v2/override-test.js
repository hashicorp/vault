/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { module, test } from 'qunit';
import { configBuilder } from 'vault/forms/v2/overrides/override-field';
import V2Form from 'vault/forms/v2/v2-form';
import { setupRenderingTest } from 'vault/tests/helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | form/v2/override', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.generatedFormConfig = {
      name: 'test-form',
      path: '/test/path',
      title: 'Test Form',
      payload: {
        username: '',
        email: '',
        enabled: false,
      },
      submit: async () => ({ success: true }),
      sections: [
        {
          name: 'section1',
          title: 'User Information',
          description: 'Enter user details',
          fields: [
            {
              name: 'username',
              type: 'TextInput',
              label: 'Username',
              helperText: 'Enter your username',
            },
            {
              name: 'email',
              type: 'TextInput',
              label: 'Email',
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

    this.form = new V2Form(this.generatedFormConfig);
  });

  test('it renders the form according to the generated config', async function (assert) {
    await render(hbs`
      <Form::V2::Renderer @form={{this.form}} @renderFields={{true}}/>
    `);

    assert.dom('form').exists('renders form element');
  });

  test('it renders an added field', async function (assert) {
    const overriddenConfig = configBuilder(this.generatedFormConfig)
      .addField('section1', {
        name: 'password',
        type: 'MaskedInput',
        label: 'Password',
      })
      .build();
    this.form = new V2Form(overriddenConfig);
    await render(hbs`
      <Form::V2::Renderer @form={{this.form}} @renderFields={{true}}/>
    `);

    assert.dom(GENERAL.maskedInput).exists('renders added field input');
    assert.dom(GENERAL.formLabel('Password')).includesText('Password', 'renders added field label');
  });

  test('it renders an updated field', async function (assert) {
    const overriddenConfig = configBuilder(this.generatedFormConfig)
      .updateField('section1', 'username', {
        type: 'TextArea',
      })
      .build();
    this.form = new V2Form(overriddenConfig);
    await render(hbs`
      <Form::V2::Renderer @form={{this.form}} @renderFields={{true}}/>
    `);

    assert.dom(GENERAL.textareaByAttr('username')).exists('renders updated TextArea field type');
    assert
      .dom(GENERAL.inputByAttr('username'))
      .doesNotExist('does not render overridden TextInput field type');
    assert.dom(GENERAL.formLabel('Username')).includesText('Username', 'renders unaltered field label');
    assert.dom(GENERAL.helpText).includesText('Enter your username', 'renders unaltered field helper text');
  });

  test('it does not render a removed field', async function (assert) {
    const overriddenConfig = configBuilder(this.generatedFormConfig)
      .removeField('section1', 'username')
      .build();
    this.form = new V2Form(overriddenConfig);
    await render(hbs`
      <Form::V2::Renderer @form={{this.form}} @renderFields={{true}}/>
    `);

    assert.dom(GENERAL.inputByAttr('username')).doesNotExist('does not render the removed username field');
    assert
      .dom(GENERAL.formLabel('Username'))
      .doesNotExist('does not render the removed username field label');
    assert.dom(GENERAL.helpText).doesNotExist('does not render the removed username field helper text');
  });

  test('it renders a field that has been reorganized into a different section', async function (assert) {
    const overriddenConfig = configBuilder(this.generatedFormConfig)
      .moveField('username', 'section1', 'section2')
      .build();
    this.form = new V2Form(overriddenConfig);
    await render(hbs`
      <Form::V2::Renderer @form={{this.form}} @renderFields={{true}}/>
    `);

    // Verify username field exists
    assert.dom(GENERAL.inputByAttr('username')).exists('username field is rendered');

    // Verify username is no longer in section1 with email
    const section1 = overriddenConfig.sections.find((s) => s.name === 'section1');
    const section2 = overriddenConfig.sections.find((s) => s.name === 'section2');

    assert.notOk(
      section1.fields.find((f) => f.name === 'username'),
      'username field removed from section1'
    );
    assert.ok(
      section2.fields.find((f) => f.name === 'username'),
      'username field added to section2'
    );

    // Verify section2 now has 2 fields (enabled + username)
    assert.strictEqual(section2.fields.length, 2, 'section2 has 2 fields after move');
  });

  test('it can move a field within a section using reorderFields', async function (assert) {
    const overriddenConfig = configBuilder(this.generatedFormConfig)
      .reorderFields('section1', ['email', 'username'])
      .build();
    this.form = new V2Form(overriddenConfig);
    await render(hbs`
      <Form::V2::Renderer @form={{this.form}} @renderFields={{true}}/>
    `);

    // Verify both fields exist
    assert.dom(GENERAL.inputByAttr('username')).exists('username field is rendered');
    assert.dom(GENERAL.inputByAttr('email')).exists('email field is rendered');

    // Verify field order in config
    const section1 = overriddenConfig.sections.find((s) => s.name === 'section1');
    assert.strictEqual(section1.fields[0].name, 'email', 'email is first field in section1');
    assert.strictEqual(section1.fields[1].name, 'username', 'username is second field in section1');

    // Verify original order was reversed
    assert.strictEqual(
      this.generatedFormConfig.sections[0].fields[0].name,
      'username',
      'original config still has username first'
    );
    assert.strictEqual(
      this.generatedFormConfig.sections[0].fields[1].name,
      'email',
      'original config still has email second'
    );
  });

  test('it can change labels, helper text, and placeholder values', async function (assert) {
    const overriddenConfig = configBuilder(this.generatedFormConfig)
      .updateField('section1', 'username', {
        label: 'Updated Label',
        helperText: 'Updated helper text',
        placeholder: 'Updated placeholder text',
      })
      .build();
    this.form = new V2Form(overriddenConfig);
    await render(hbs`
      <Form::V2::Renderer @form={{this.form}} @renderFields={{true}}/>
    `);

    assert
      .dom(GENERAL.formLabel('Updated Label'))
      .includesText('Updated Label', 'renders updated username label');
    assert.dom(GENERAL.helpText).includesText('Updated helper text', 'renders updated helper text');
    assert
      .dom(GENERAL.inputByAttr('username'))
      .hasAttribute('placeholder', 'Updated placeholder text', 'renders updated placeholder text');
  });

  test('it can add field validation', async function (assert) {
    const overriddenConfig = configBuilder(this.generatedFormConfig)
      .updateField('section1', 'username', {
        validations: [
          { type: 'required', message: 'Username is required' },
          {
            type: 'minLength',
            options: { minLength: 3 },
            message: 'Username must be at least 3 characters',
          },
          {
            type: 'maxLength',
            options: { maxLength: 20 },
            message: 'Username must be at most 20 characters',
          },
        ],
        isRequired: true,
      })
      .build();
    this.form = new V2Form(overriddenConfig);
    await render(hbs`
      <Form::V2::Renderer @form={{this.form}} @renderFields={{true}}/>
    `);

    const section1 = overriddenConfig.sections.find((s) => s.name === 'section1');
    const usernameField = section1.fields.find((f) => f.name === 'username');

    assert.ok(usernameField.isRequired, 'username field is marked as required');
    assert.strictEqual(usernameField.validations.length, 3, 'username field has 3 validations');
    assert.strictEqual(usernameField.validations[0].type, 'required', 'first validation is required');
    assert.strictEqual(usernameField.validations[1].type, 'minLength', 'second validation is minLength');
    assert.strictEqual(usernameField.validations[2].type, 'maxLength', 'third validation is maxLength');
  });
});
