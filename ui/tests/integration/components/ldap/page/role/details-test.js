/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { duration } from 'core/helpers/format-duration';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | ldap | Page::Role::Details', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.backend = 'ldap-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);
    this.model = {
      capabilities: {
        canDelete: true,
        canEdit: true,
        canReadCreds: true,
        canRotateStaticCreds: true,
      },
    };
    this.renderComponent = (type) => {
      this.model.role = this.server.create('ldap-role', type);
      this.breadcrumbs = [
        { label: this.backend, route: 'overview' },
        { label: 'Roles', route: 'roles' },
        { label: this.model.role.name },
      ];
      return render(hbs`<Page::Role::Details @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`, {
        owner: this.engine,
      });
    };
  });

  test('it should render header with role name and breadcrumbs', async function (assert) {
    await this.renderComponent('static');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText(this.model.role.name, 'Role name renders in header');
    assert
      .dom('[data-test-breadcrumbs] li:nth-child(1)')
      .containsText(this.backend, 'Overview breadcrumb renders');
    assert.dom('[data-test-breadcrumbs] li:nth-child(2) a').containsText('Roles', 'Roles breadcrumb renders');
    assert
      .dom('[data-test-breadcrumbs] li:nth-child(3)')
      .containsText(this.model.role.name, 'Role breadcrumb renders');
  });

  test('it should render page header dropdown actions', async function (assert) {
    assert.expect(7);

    await this.renderComponent('static');
    await click(GENERAL.dropdownToggle('Manage'));

    assert.dom(GENERAL.menuItem('Delete role')).hasText('Delete role', 'Delete action renders');
    assert
      .dom(GENERAL.button('Get credentials'))
      .hasText('Get credentials', 'Get credentials action renders');
    assert
      .dom(GENERAL.menuItem('Rotate credentials'))
      .exists('Rotate credentials action renders for static role');
    assert.dom(GENERAL.menuItem('Edit role')).hasText('Edit role', 'Edit action renders');

    this.model.capabilities.canRotateStaticCreds = false;
    await this.renderComponent('dynamic');
    // defined after render so this.model is defined
    const deleteStub = sinon
      .stub(this.owner.lookup('service:api').secrets, 'ldapDeleteDynamicRole')
      .resolves();
    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    assert
      .dom('[data-test-rotate-credentials]')
      .doesNotExist('Rotate credentials action is hidden for dynamic role');
    await click(GENERAL.dropdownToggle('Manage'));
    await click(GENERAL.menuItem('Delete role'));
    await click(GENERAL.confirmButton);

    assert.true(
      deleteStub.calledWith(this.model.role.name, this.backend),
      'Delete API called with correct parameters'
    );
    assert.true(
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
        const modelValue = this.model.role[field.key];
        const isDuration = ['TTL', 'Max TTL', 'Rotation period'].includes(field.label);
        const value = isDuration ? duration([modelValue]) : modelValue;
        assert.dom(`[data-test-row-value="${field.label}"]`).hasText(value, `${field.label} value renders`);
      });
    }
  });

  test('it should show password row when canReadCreds is true and hide it when false', async function (assert) {
    assert.expect(2);

    // canReadCreds: true (set in beforeEach) — password row should be visible
    await this.renderComponent('static');
    assert.dom('[data-test-row-label="Password"]').exists('Password row renders when canReadCreds is true');

    // canReadCreds: false — password row should be hidden
    this.model.capabilities.canReadCreds = false;
    await this.renderComponent('static');
    assert
      .dom('[data-test-row-label="Password"]')
      .doesNotExist('Password row is hidden when canReadCreds is false');
  });

  // The credential is fetched only when the user asks to see it, so simply opening the details
  // page never pulls a secret the user did not request.
  test('it should fetch the password only on the first reveal', async function (assert) {
    const credsStub = sinon
      .stub(this.owner.lookup('service:api').secrets, 'ldapRequestStaticRoleCredentials')
      .resolves({ data: { password: 'super-secret' } });

    await this.renderComponent('static');
    assert.false(credsStub.called, 'no credential request is made on render');

    await click(GENERAL.button('toggle-masked'));
    assert.true(credsStub.calledOnce, 'the credential is fetched when the value is revealed');
    assert.true(
      credsStub.calledWith(this.model.role.name, this.backend),
      'the request targets the role on the current mount'
    );
    assert.dom(GENERAL.maskedInput).hasText('super-secret', 'the password is displayed');

    // re-mask and reveal again — the cached value is reused
    await click(GENERAL.button('toggle-masked'));
    await click(GENERAL.button('toggle-masked'));
    assert.true(credsStub.calledOnce, 'the cached password is reused rather than refetched');
  });

  // The edit form seeds itself from this same role model, so a password stored there would
  // pre-fill the edit page's password field with the real secret.
  test('it should not write the password onto the shared route model', async function (assert) {
    sinon
      .stub(this.owner.lookup('service:api').secrets, 'ldapRequestStaticRoleCredentials')
      .resolves({ data: { password: 'super-secret' } });

    await this.renderComponent('static');
    await click(GENERAL.button('toggle-masked'));

    assert.notOk(this.model.role.password, 'the role model carries no plaintext password');
  });

  test('it should render the minus icon when the role has no password', async function (assert) {
    sinon
      .stub(this.owner.lookup('service:api').secrets, 'ldapRequestStaticRoleCredentials')
      .resolves({ data: { password: '' } });

    await this.renderComponent('static');
    await click(GENERAL.button('toggle-masked'));

    assert
      .dom(`${GENERAL.maskedInput} .hds-icon-minus`)
      .exists('an empty password renders as a minus icon rather than a blank row');
  });

  test('it should render an inline error when the credential request fails', async function (assert) {
    sinon
      .stub(this.owner.lookup('service:api').secrets, 'ldapRequestStaticRoleCredentials')
      .rejects({ status: 403 });

    await this.renderComponent('static');
    await click(GENERAL.button('toggle-masked'));

    assert.dom(GENERAL.inlineAlert).exists('an inline alert renders when the request fails');
  });
});
