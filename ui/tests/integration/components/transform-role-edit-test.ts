/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import RoleForm from 'vault/forms/transform/role';
import sinon from 'sinon';
import type ApiService from 'vault/services/api';

module('Integration | Component | transform-role-edit', function (hooks) {
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
      new RoleForm(
        {
          name: 'my-role',
          transformations: ['my-transformation'],
          backend: 'transform',
        },
        { isNew: false }
      )
    );
    this.set('mode', 'show');

    await render(
      hbs`<TransformRoleEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />`
    );

    assert.dom('[data-test-edit-link]').exists('renders toolbar edit link');
    assert.dom('[data-test-field]').doesNotExist('does not render form fields in show mode');
  });

  test('it renders in create mode', async function (assert) {
    this.set('form', new RoleForm({ backend: 'transform' }, { isNew: true }));
    this.set('mode', 'create');

    await render(
      hbs`<TransformRoleEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />`
    );

    assert.dom('[data-test-submit]').exists('renders submit button');
    assert.dom('[data-test-submit]').hasText('Create role');
  });

  test('it renders in edit mode', async function (assert) {
    this.set(
      'form',
      new RoleForm(
        {
          name: 'my-role',
          transformations: ['my-transformation'],
          backend: 'transform',
        },
        { isNew: false }
      )
    );
    this.set('mode', 'edit');

    await render(
      hbs`<TransformRoleEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />`
    );

    assert.dom('[data-test-submit]').exists('renders submit button');
    assert.dom('[data-test-submit]').hasText('Save');
    assert.dom('[data-test-input="name"]').hasAttribute('readonly', '', 'name is readonly in edit mode');
  });

  test('it calls onDelete and transitions to list', async function (assert) {
    const api = this.owner.lookup('service:api') as unknown as ApiService;
    const deleteStub = sinon.stub(api.secrets, 'transformDeleteRole').resolves();

    this.set('form', new RoleForm({ name: 'my-role', backend: 'transform' }, { isNew: false }));
    this.set('mode', 'show');

    await render(
      hbs`<TransformRoleEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />`
    );

    await click('[data-test-delete]');

    assert.ok(deleteStub.calledWith('my-role', 'transform'), 'calls transformDeleteRole with correct args');
    deleteStub.restore();
  });
});
