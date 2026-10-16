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

  // An empty dropdown would open a blank popover, so the toggle is omitted entirely.
  test('it should hide the Manage dropdown when no actions are permitted', async function (assert) {
    this.model.capabilities = {
      canDelete: false,
      canEdit: false,
      canReadCreds: false,
      canRotateStaticCreds: false,
    };

    for (const type of ['static', 'dynamic']) {
      await this.renderComponent(type);
      assert.dom(GENERAL.dropdownToggle('Manage')).doesNotExist(`Manage is hidden for a ${type} role`);
    }
  });

  test('it should show the Manage dropdown with only the permitted action', async function (assert) {
    const actions = {
      canEdit: 'Edit role',
      canDelete: 'Delete role',
      canRotateStaticCreds: 'Rotate credentials',
    };

    for (const [capability, label] of Object.entries(actions)) {
      this.model.capabilities = {
        canDelete: false,
        canEdit: false,
        canReadCreds: false,
        canRotateStaticCreds: false,
        [capability]: true,
      };
      await this.renderComponent('static');
      await click(GENERAL.dropdownToggle('Manage'));

      assert.dom(GENERAL.menuItem()).exists({ count: 1 }, `only one action renders for ${capability}`);
      assert.dom(GENERAL.menuItem(label)).exists(`${label} renders for ${capability}`);
    }
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

  module('rotating credentials', function (hooks) {
    hooks.beforeEach(function () {
      const { secrets } = this.owner.lookup('service:api');
      this.rotateStub = sinon.stub(secrets, 'ldapRotateStaticRole').resolves();
      this.credsStub = sinon.stub(secrets, 'ldapRequestStaticRoleCredentials');
      this.credsStub.onFirstCall().resolves({ data: { password: 'old-password' } });
      this.credsStub.onSecondCall().resolves({ data: { password: 'new-password' } });
      this.rotate = async () => {
        await click(GENERAL.dropdownToggle('Manage'));
        await click(GENERAL.menuItem('Rotate credentials'));
        await click(GENERAL.confirmButton);
      };
    });

    // Rotation invalidates the revealed password, so the row must show the new one rather than
    // going blank and making the user click reveal again.
    test('it should refresh a revealed password after rotating', async function (assert) {
      await this.renderComponent('static');
      await click(GENERAL.button('toggle-masked'));
      assert.dom(GENERAL.maskedInput).hasText('old-password', 'the current password is revealed');

      await this.rotate();

      assert.true(this.rotateStub.calledOnce, 'the role is rotated');
      assert.true(this.credsStub.calledTwice, 'the password is fetched again after rotating');
      assert.dom(GENERAL.maskedInput).hasText('new-password', 'the rotated password is displayed');
    });

    test('it should copy the rotated password after rotating', async function (assert) {
      const clipboardStub = sinon.stub(navigator.clipboard, 'writeText').resolves();

      await this.renderComponent('static');
      await click(GENERAL.button('toggle-masked'));
      await this.rotate();
      await click(GENERAL.copyButton);

      assert.true(clipboardStub.calledOnceWith('new-password'), 'the rotated password is copied');
    });

    // Opening the details page or rotating must not pull a secret the user never asked to see.
    test('it should not fetch the password after rotating when it was never revealed', async function (assert) {
      await this.renderComponent('static');
      await this.rotate();

      assert.true(this.rotateStub.calledOnce, 'the role is rotated');
      assert.false(this.credsStub.called, 'no credential request is made');
      assert.dom(GENERAL.maskedInput).hasText('***********', 'the password stays masked');
    });

    test('it should render an inline error when the password cannot be refreshed after rotating', async function (assert) {
      const flashSuccessSpy = sinon.spy(this.owner.lookup('service:flash-messages'), 'success');
      this.credsStub.onSecondCall().rejects({ status: 403 });

      await this.renderComponent('static');
      await click(GENERAL.button('toggle-masked'));
      await this.rotate();

      assert.true(
        flashSuccessSpy.calledWith('Credentials successfully rotated.'),
        'the rotation still reports success'
      );
      assert.dom(GENERAL.inlineAlert).exists('an inline alert explains why the password is missing');
      assert
        .dom(GENERAL.maskedInput)
        .doesNotIncludeText('old-password', 'the invalidated password is no longer displayed');
    });

    test('it should fetch the password after rotating when the first reveal failed', async function (assert) {
      this.credsStub.onFirstCall().rejects({ status: 500 });

      await this.renderComponent('static');
      await click(GENERAL.button('toggle-masked'));
      assert.dom(GENERAL.inlineAlert).exists('the failed reveal shows an inline error');

      await this.rotate();

      assert.true(this.credsStub.calledTwice, 'the password is fetched again after rotating');
      assert.dom(GENERAL.maskedInput).hasText('new-password', 'the rotated password is displayed');
      assert.dom(GENERAL.inlineAlert).doesNotExist('the earlier error is cleared');
    });
  });
});
