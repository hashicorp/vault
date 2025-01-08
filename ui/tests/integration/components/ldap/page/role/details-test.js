/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { duration } from 'core/helpers/format-duration';
import { ldapRoleID } from 'vault/adapters/ldap/role';

module('Integration | Component | ldap | Page::Role::Details', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.server.post('/sys/capabilities-self', () => ({
      data: {
        capabilities: ['root'],
      },
    }));
    this.renderComponent = (type) => {
      const data = this.server.create('ldap-role', type);
      data.id = ldapRoleID(type, data.name);
      const store = this.owner.lookup('service:store');
      store.pushPayload('ldap/role', {
        modelName: 'ldap/role',
        backend: 'ldap-test',
        type,
        ...data,
      });
      this.model = store.peekRecord('ldap/role', ldapRoleID(type, data.name));
      this.breadcrumbs = [
        { label: this.model.backend, route: 'overview' },
        { label: 'Roles', route: 'roles' },
        { label: this.model.name },
      ];
      return render(hbs`<Page::Role::Details @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`, {
        owner: this.engine,
      });
    };
  });

  test('it should render header with role name and breadcrumbs', async function (assert) {
    await this.renderComponent('static');
    assert.dom('[data-test-header-title]').hasText(this.model.name, 'Role name renders in header');
    assert
      .dom('[data-test-breadcrumbs] li:nth-child(1)')
      .containsText(this.model.backend, 'Overview breadcrumb renders');
    assert.dom('[data-test-breadcrumbs] li:nth-child(2) a').containsText('Roles', 'Roles breadcrumb renders');
    assert
      .dom('[data-test-breadcrumbs] li:nth-child(3)')
      .containsText(this.model.name, 'Role breadcrumb renders');
  });

  test('it should render toolbar actions', async function (assert) {
    assert.expect(7);

    await this.renderComponent('static');

    assert.dom('[data-test-delete]').hasText('Delete role', 'Delete action renders');
    assert.dom('[data-test-get-credentials]').hasText('Get credentials', 'Get credentials action renders');
    assert.dom('[data-test-rotate-credentials]').exists('Rotate credentials action renders for static role');
    assert.dom('[data-test-edit]').hasText('Edit role', 'Edit action renders');

    await this.renderComponent('dynamic');
    // defined after render so this.model is defined
    this.server.delete(`/${this.model.backend}/role/${this.model.name}`, () => {
      assert.ok(true, 'Request made to delete role');
      return;
    });
    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    assert
      .dom('[data-test-rotate-credentials]')
      .doesNotExist('Rotate credentials action is hidden for dynamic role');

    await click('[data-test-delete]');
    await click('[data-test-confirm-button]');
    assert.ok(
      transitionStub.calledWith('vault.cluster.secrets.backend.ldap.roles'),
      'Transitions to roles route on delete success'
    );
  });

  test('it should render details fields', async function (assert) {
    assert.expect(26);

    const fields = [
      { label: 'Role name', key: 'name' },
      { label: 'Role type', key: 'type' },
      { label: 'Distinguished name', key: 'dn', type: 'static' },
      { label: 'Username', key: 'username', type: 'static' },
      { label: 'Rotation period', key: 'rotation_period', type: 'static' },
      { label: 'TTL', key: 'default_ttl', type: 'dynamic' },
      { label: 'Max TTL', key: 'max_ttl', type: 'dynamic' },
      { label: 'Username template', key: 'username_template', type: 'dynamic' },
      { label: 'Creation LDIF', key: 'creation_ldif', type: 'dynamic' },
      { label: 'Deletion LDIF', key: 'deletion_ldif', type: 'dynamic' },
      { label: 'Rollback LDIF', key: 'rollback_ldif', type: 'dynamic' },
    ];

    for (const type of ['static', 'dynamic']) {
      await this.renderComponent(type);

      const typeFields = fields.filter((field) => !field.type || field.type === type);
      typeFields.forEach((field) => {
        assert
          .dom(`[data-test-row-label="${field.label}"]`)
          .hasText(field.label, `${field.label} label renders`);
        const modelValue = this.model[field.key];
        const isDuration = ['TTL', 'Max TTL', 'Rotation period'].includes(field.label);
        const value = isDuration ? duration([modelValue]) : modelValue;
        assert.dom(`[data-test-row-value="${field.label}"]`).hasText(value, `${field.label} value renders`);
      });
    }
  });
});
