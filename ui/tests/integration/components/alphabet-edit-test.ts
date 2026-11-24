/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import AlphabetForm from 'vault/forms/transform/alphabet';
import sinon from 'sinon';
import type ApiService from 'vault/services/api';

module('Integration | Component | alphabet-edit', function (hooks) {
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
      new AlphabetForm(
        {
          name: 'my-alphabet',
          alphabet: 'abcdefghijklmnopqrstuvwxyz',
          backend: 'transform',
        },
        { isNew: false }
      )
    );
    this.set('mode', 'show');

    await render(
      hbs`<AlphabetEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />`
    );

    assert.dom('[data-test-edit-link]').exists('renders toolbar edit link');
    assert.dom('[data-test-field]').doesNotExist('does not render form fields in show mode');
  });

  test('it renders in create mode', async function (assert) {
    this.set('form', new AlphabetForm({ backend: 'transform' }, { isNew: true }));
    this.set('mode', 'create');

    await render(
      hbs`<AlphabetEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />`
    );

    assert.dom('[data-test-submit]').exists('renders submit button');
    assert.dom('[data-test-submit]').hasText('Create alphabet');
    const fields = findAll('[data-test-field]');
    assert.strictEqual(fields.length, 2, 'renders name and alphabet fields');
  });

  test('it renders in edit mode', async function (assert) {
    this.set(
      'form',
      new AlphabetForm(
        {
          name: 'my-alphabet',
          alphabet: 'abcdefghijklmnopqrstuvwxyz',
          backend: 'transform',
        },
        { isNew: false }
      )
    );
    this.set('mode', 'edit');

    await render(
      hbs`<AlphabetEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />`
    );

    assert.dom('[data-test-submit]').exists('renders submit button');
    assert.dom('[data-test-submit]').hasText('Save');
    assert.dom('[data-test-input="name"]').hasAttribute('readonly', '', 'name is readonly in edit mode');
  });

  test('it calls onDelete and transitions to list', async function (assert) {
    const api = this.owner.lookup('service:api') as unknown as ApiService;
    const deleteStub = sinon.stub(api.secrets, 'transformDeleteAlphabet').resolves();

    this.set('form', new AlphabetForm({ name: 'my-alphabet', backend: 'transform' }, { isNew: false }));
    this.set('mode', 'show');

    await render(
      hbs`<AlphabetEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />`
    );

    await click('[data-test-delete]');

    assert.ok(
      deleteStub.calledWith('my-alphabet', 'transform'),
      'calls transformDeleteAlphabet with correct args'
    );
    deleteStub.restore();
  });
});
