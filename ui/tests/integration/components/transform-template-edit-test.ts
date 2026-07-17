/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import TemplateForm from 'vault/forms/transform/template';
import sinon from 'sinon';
import type ApiService from 'vault/services/api';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | transform-template-edit', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    const router = this.owner.lookup('service:router') as unknown as Record<string, unknown>;
    router['transitionTo'] = sinon.stub();

    this.set('capabilities', {
      canDelete: true,
      canUpdate: true,
      canRead: true,
    });
  });

  test('it renders in show mode', async function (assert) {
    this.set(
      'form',
      new TemplateForm(
        {
          name: 'my-template',
          pattern: '[0-9]{9}',
          alphabet: ['builtin/numerics'],
          encode_format: 'local',
          decode_formats: { local: '[0-9]{9}' },
          backend: 'transform',
        },
        { isNew: false }
      )
    );
    this.set('mode', 'show');

    await render(
      hbs`<TransformTemplateEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />`
    );

    assert.dom('[data-test-edit-link]').exists('renders toolbar edit link');
    assert.dom('[data-test-field]').doesNotExist('does not render form fields in show mode');
  });

  test('it renders in create mode', async function (assert) {
    this.set('form', new TemplateForm({ backend: 'transform' }, { isNew: true }));
    this.set('mode', 'create');

    await render(
      hbs`<TransformTemplateEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />`
    );

    assert.dom(GENERAL.submitButton).exists('renders submit button');
    assert.dom(GENERAL.submitButton).hasText('Create template');
    const fields = findAll('[data-test-field]');
    // name, pattern, alphabet (encodeFormat/decodeFormats skipped — handled by TransformAdvancedTemplating)
    assert.strictEqual(fields.length, 3, 'renders 3 form fields for create mode');
  });

  test('it renders in edit mode', async function (assert) {
    this.set(
      'form',
      new TemplateForm(
        {
          name: 'my-template',
          pattern: '[0-9]{9}',
          alphabet: ['builtin/numerics'],
          backend: 'transform',
        },
        { isNew: false }
      )
    );
    this.set('mode', 'edit');

    await render(
      hbs`<TransformTemplateEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />`
    );

    assert.dom(GENERAL.submitButton).exists('renders submit button');
    assert.dom(GENERAL.submitButton).hasText('Save');
  });

  test('it calls onDelete and transitions to list', async function (assert) {
    const api = this.owner.lookup('service:api') as unknown as ApiService;
    const deleteStub = sinon.stub(api.secrets, 'transformDeleteTemplate').resolves();

    this.set('form', new TemplateForm({ name: 'my-template', backend: 'transform' }, { isNew: false }));
    this.set('mode', 'show');

    await render(
      hbs`<TransformTemplateEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />`
    );

    await click('[data-test-delete]');

    assert.ok(
      deleteStub.calledWith('my-template', 'transform'),
      'calls transformDeleteTemplate with correct args'
    );
    deleteStub.restore();
  });
});
