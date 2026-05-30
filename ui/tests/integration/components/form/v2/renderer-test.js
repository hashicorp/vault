/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { render, waitFor } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { module, test } from 'qunit';
import V2Form from 'vault/forms/v2/v2-form';
import { setupRenderingTest } from 'vault/tests/helpers';

module('Integration | Component | form/v2/renderer', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.formConfig = {
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

    this.form = new V2Form(this.formConfig);
  });

  test('it renders the form element', async function (assert) {
    await render(hbs`
      <Form::V2::Renderer @form={{this.form}} />
    `);

    assert.dom('form').exists('renders form element');
  });

  test('it renders sections and fields when renderFields is true', async function (assert) {
    await render(hbs`
      <Form::V2::Renderer @form={{this.form}} @renderFields={{true}} />
    `);

    await waitFor('h2');
    assert.dom('h2').exists({ count: 2 }, 'renders section headers');
    assert.dom('label').includesText('Username', 'renders field content');
  });

  test('it renders an error alert when error is provided', async function (assert) {
    this.error = 'Something went wrong';

    await render(hbs`
      <Form::V2::Renderer @form={{this.form}} @error={{this.error}} @renderFields={{true}} />
    `);

    assert.dom('.hds-alert').exists('renders error alert');
    assert.dom('.hds-alert__description').hasText('Something went wrong', 'renders error message');
  });

  test('it yields custom content', async function (assert) {
    await render(hbs`
      <Form::V2::Renderer @form={{this.form}} @renderFields={{true}} as |Form|>
        <Form.Section>
          <div data-test-custom-content>Custom submit button</div>
        </Form.Section>
      </Form::V2::Renderer>
    `);

    assert
      .dom('[data-test-custom-content]')
      .hasText('Custom submit button', 'renders yielded custom content');
  });

  test('it does not render fields when renderFields is false', async function (assert) {
    await render(hbs`
      <Form::V2::Renderer @form={{this.form}} @renderFields={{false}} />
    `);

    assert.dom('h2').doesNotExist('does not render section headers');
    assert.dom('label').doesNotExist('does not render fields');
  });

  test('it handles forms with empty sections', async function (assert) {
    this.formWithEmptySection = new V2Form({
      name: 'form-empty-section',
      path: '/test',
      payload: {},
      submit: async () => ({}),
      sections: [
        {
          name: 'empty-section',
          title: 'Empty Section',
          fields: [],
        },
      ],
    });

    await render(hbs`
      <Form::V2::Renderer @form={{this.formWithEmptySection}} @renderFields={{true}} />
    `);

    assert.dom('h2').hasText('Empty Section', 'renders section title');
    assert.dom('label').doesNotExist('does not render field content');
  });

  test('it passes attributes to the form element', async function (assert) {
    await render(hbs`
      <Form::V2::Renderer
        @form={{this.form}}
        @renderFields={{true}}
        data-test-custom-form
        class="custom-class"
      />
    `);

    assert.dom('form').hasAttribute('data-test-custom-form', '', 'passes data attributes');
    assert.dom('form').hasClass('custom-class', 'passes CSS classes');
  });
});
