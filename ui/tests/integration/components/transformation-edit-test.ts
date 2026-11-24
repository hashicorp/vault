/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click, fillIn, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import TransformationForm from 'vault/forms/transform/transformation';
import sinon from 'sinon';
import type ApiService from 'vault/services/api';

module('Integration | Component | transformation-edit', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    const router = this.owner.lookup('service:router') as unknown as Record<string, unknown>;
    router['transitionTo'] = sinon.stub();

    this.set('capabilities', {
      canDelete: true,
      canUpdate: true,
      canRead: true,
    });

    // Stub list fetches called in constructor to avoid real API calls
    const api = this.owner.lookup('service:api') as unknown as ApiService;
    sinon.stub(api.secrets, 'transformListRoles').resolves({ keys: [] });
    sinon.stub(api.secrets, 'transformListTemplates').resolves({ keys: [] });
  });

  hooks.afterEach(function () {
    sinon.restore();
  });

  test('it renders in show mode', async function (assert) {
    this.set(
      'form',
      new TransformationForm(
        {
          name: 'my-transformation',
          type: 'fpe',
          allowed_roles: [],
          backend: 'transform',
        },
        { isNew: false }
      )
    );
    this.set('mode', 'show');

    await render(
      hbs`<TransformationEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />`
    );

    assert.dom('[data-test-edit-link]').exists('renders toolbar edit link');
    assert.dom('[data-test-field]').doesNotExist('does not render form fields in show mode');
  });

  test('it renders in create mode with fpe type', async function (assert) {
    this.set('form', new TransformationForm({ backend: 'transform', type: 'fpe' }, { isNew: true }));
    this.set('mode', 'create');

    await render(
      hbs`<TransformationEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />`
    );

    assert.dom('[data-test-submit]').exists('renders submit button');
    assert.dom('[data-test-submit]').hasText('Create transformation');
    // fpe: name, type, deletion_allowed, tweak_source (FormField) + template, allowed_roles (SearchSelect with data-test-field)
    const fields = findAll('[data-test-field]');
    assert.strictEqual(fields.length, 6, 'renders 6 fields for fpe type (4 FormField + 2 SearchSelect)');
  });

  test('it renders masking-specific field for masking type', async function (assert) {
    this.set('form', new TransformationForm({ backend: 'transform', type: 'masking' }, { isNew: true }));
    this.set('mode', 'create');

    await render(
      hbs`<TransformationEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />`
    );

    // masking: name, type, deletion_allowed, masking_character (FormField) + template, allowed_roles (SearchSelect with data-test-field)
    const fields = findAll('[data-test-field]');
    assert.strictEqual(fields.length, 6, 'renders 6 fields for masking type (4 FormField + 2 SearchSelect)');
    assert.dom('[data-test-input="masking_character"]').exists('renders masking_character field');
    assert
      .dom('[data-test-input="tweak_source"]')
      .doesNotExist('does not render tweak_source for masking type');
  });

  test('it renders tokenization-specific fields for tokenization type', async function (assert) {
    this.set('form', new TransformationForm({ backend: 'transform', type: 'tokenization' }, { isNew: true }));
    this.set('mode', 'create');

    await render(
      hbs`<TransformationEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />`
    );

    // tokenization shows: name, type, deletion_allowed, mapping_mode, convergent, max_ttl, stores, allowed_roles (SearchSelect)
    assert.dom('[data-test-input="mapping_mode"]').exists('renders mapping_mode field');
    assert.dom('[data-test-input="convergent"]').exists('renders convergent field');
    assert.dom('[data-test-input="max_ttl"]').exists('renders max_ttl field');
    assert
      .dom('[data-test-input="tweak_source"]')
      .doesNotExist('does not render tweak_source for tokenization type');
    assert
      .dom('[data-test-input="masking_character"]')
      .doesNotExist('does not render masking_character for tokenization type');
  });

  test('it renders in edit mode with name readonly', async function (assert) {
    this.set(
      'form',
      new TransformationForm(
        {
          name: 'my-transformation',
          type: 'fpe',
          allowed_roles: [],
          backend: 'transform',
        },
        { isNew: false }
      )
    );
    this.set('mode', 'edit');

    await render(
      hbs`<TransformationEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />`
    );

    assert.dom('[data-test-submit]').hasText('Save');
    assert.dom('[data-test-input="name"]').hasAttribute('readonly', '', 'name is readonly in edit mode');
  });

  test('it calls onDelete and transitions to list', async function (assert) {
    const api = this.owner.lookup('service:api') as unknown as ApiService;
    const deleteStub = sinon.stub(api.secrets, 'transformDeleteTransformation').resolves();

    this.set(
      'form',
      new TransformationForm(
        { name: 'my-transformation', type: 'fpe', allowed_roles: [], backend: 'transform' },
        { isNew: false }
      )
    );
    this.set('mode', 'show');

    await render(
      hbs`<TransformationEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />`
    );

    await click('[data-test-delete]');
    await fillIn('[data-test-confirmation-modal-input="Delete transformation"]', 'my-transformation');
    await click('[data-test-confirm-button="Delete transformation"]');

    assert.ok(
      deleteStub.calledWith('my-transformation', 'transform'),
      'calls transformDeleteTransformation with correct args'
    );
  });

  test('it shows edit warning modal when transformation has allowed roles', async function (assert) {
    this.set(
      'form',
      new TransformationForm(
        {
          name: 'my-transformation',
          type: 'fpe',
          allowed_roles: ['my-role'],
          backend: 'transform',
        },
        { isNew: false }
      )
    );
    this.set('mode', 'show');

    await render(
      hbs`<TransformationEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />`
    );

    await click('[data-test-edit-link]');

    assert.dom('#transformation-edit-modal').exists('shows edit warning modal when transformation has roles');
    assert.dom('[data-test-edit-confirm-button]').exists('renders confirm button in modal');
  });
});
