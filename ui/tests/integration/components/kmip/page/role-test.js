/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { click, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { capitalize } from '@ember/string';
import operationGroups from 'kmip/helpers/operation-groups';

module('Integration | Component | kmip | Page::Role', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kmip');

  hooks.beforeEach(function () {
    this.backend = 'kmip-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    this.roleName = 'role-1';
    this.scopeName = 'scope-1';
    this.capabilities = { canDelete: true, canUpdate: true };

    this.getRole = (allOrNone) => {
      const tlsOptions = {
        tls_client_key_bits: 521,
        tls_client_key_type: 'ec',
        tls_client_ttl: 86400,
      };
      let operations = {
        operation_activate: true,
        operation_add_attribute: true,
        operation_decrypt: true,
        operation_discover_versions: true,
        operation_import: true,
        operation_locate: true,
        operation_register: true,
        operation_revoke: true,
      };
      if (allOrNone === 'all') {
        operations = { operation_all: true };
      } else if (allOrNone === 'none') {
        operations = { operation_none: true };
      }
      return { ...tlsOptions, ...operations };
    };
    this.role = this.getRole();

    // get all keys that are rendered in the operation groups
    this.operationKeys = Object.values(operationGroups()).flat();

    this.apiStub = sinon.stub(this.owner.lookup('service:api').secrets, 'kmipDeleteRole').resolves();
    this.flashStub = sinon.stub(this.owner.lookup('service:flashMessages'), 'success');
    this.routerStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    this.operationState = (assert, operation, isEnabled) => {
      const iconClass = isEnabled ? 'hds-foreground-success' : 'hds-foreground-faint';
      const iconName = isEnabled ? 'check-circle' : 'x-square';
      const label = operation.replace('operation_', '').split('_').map(capitalize).join(' ');
      const state = isEnabled ? 'enabled' : 'disabled';

      const selector = `[data-test-operation-field="${operation}"]`;
      assert
        .dom(`${selector} svg`)
        .hasClass(iconClass, `${operation} has correct icon class for ${state} state`);
      assert
        .dom(`${selector} svg`)
        .hasAttribute('data-test-icon', iconName, `${operation} has correct icon for ${state} state`);
      assert.dom(selector).containsText(label, `${operation} has correct label`);
    };

    this.renderComponent = () =>
      render(
        hbs`<Page::Role @roleName={{this.roleName}} @scopeName={{this.scopeName}} @role={{this.role}} @capabilities={{this.capabilities}} />`,
        { owner: this.engine }
      );
  });

  test('it should render/hide toolbar actions based on capabilities', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.confirmTrigger).hasText('Delete role', 'Delete role action renders in toolbar');
    assert.dom('[data-test-kmip-link-edit-role]').hasText('Edit role', 'Edit role action renders in toolbar');

    this.capabilities = { canDelete: false, canUpdate: false };
    await this.renderComponent();

    assert.dom(GENERAL.confirmTrigger).doesNotExist('Delete role action is hidden without capability');
    assert
      .dom('[data-test-kmip-link-edit-role]')
      .doesNotExist('Edit role action is hidden without capability');
  });

  test('it should delete role', async function (assert) {
    await this.renderComponent();

    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);

    assert.true(
      this.apiStub.calledWith(this.roleName, this.scopeName, this.backend),
      'API called to delete role'
    );
    assert.true(
      this.flashStub.calledWith(`Successfully deleted role ${this.roleName}`),
      'Success flash message shown'
    );
    assert.true(
      this.routerStub.calledWith('vault.cluster.secrets.backend.kmip.scope.roles', this.scopeName),
      'Transitions to roles list on delete success'
    );
  });

  test('it should render tls fields', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.infoRowValue('TLS client key bits')).hasText('521', 'TLS client key bits renders');
    assert.dom(GENERAL.infoRowValue('TLS client key type')).hasText('ec', 'TLS client key type renders');
    assert.dom(GENERAL.infoRowValue('TLS client TTL')).hasText('1 day', 'TLS client TTL renders');
  });

  test('it should render operation groups', async function (assert) {
    await this.renderComponent();

    assert
      .dom(GENERAL.infoRowLabel('Managed Cryptographic Objects'))
      .exists('Cypto operations group renders');
    assert.dom(GENERAL.infoRowLabel('Object Attributes')).exists('Attributes operations group renders');
    assert.dom(GENERAL.infoRowLabel('Server')).exists('Server operations group renders');
    assert.dom(GENERAL.infoRowLabel('Other')).exists('Other operations group renders');
  });

  test('it should mark all operations as enabled when operations_all was selected', async function (assert) {
    assert.expect(this.operationKeys.length * 3 + 1);

    this.role = this.getRole('all');
    await this.renderComponent();

    assert
      .dom(GENERAL.inlineError)
      .hasText('This role allows all KMIP operations', 'All operations enabled message renders');

    this.operationKeys.forEach((operation) => {
      this.operationState(assert, operation, true);
    });
  });

  test('it should mark operations as disabled when operations_none was selected', async function (assert) {
    assert.expect(this.operationKeys.length * 3 + 1);

    this.role = this.getRole('none');
    await this.renderComponent();

    assert
      .dom(GENERAL.inlineError)
      .doesNotExist('All operations enabled message does not render when operations are disabled');

    this.operationKeys.forEach((operation) => {
      this.operationState(assert, operation, false);
    });
  });

  test('it should correctly mark operations as enabled or disabled based on selections', async function (assert) {
    assert.expect(this.operationKeys.length * 3 + 1);

    await this.renderComponent();

    assert
      .dom(GENERAL.inlineError)
      .doesNotExist('All operations enabled message does not render when some operations are disabled');

    this.operationKeys.forEach((operation) => {
      this.operationState(assert, operation, this.role[operation] === true);
    });
  });
});
