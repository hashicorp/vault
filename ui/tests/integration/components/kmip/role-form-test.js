/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { click, fillIn, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import KmipRoleForm from 'vault/forms/secrets/kmip/role';
import operationGroups from 'kmip/helpers/operation-groups';

module('Integration | Component | kmip | RoleForm', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kmip');

  hooks.beforeEach(function () {
    this.backend = 'kmip-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    this.roleName = '';
    this.scopeName = 'scope-1';

    this.tlsOptions = {
      tls_client_key_bits: 521,
      tls_client_key_type: 'ec',
      tls_client_ttl: '86400',
    };

    this.setForm = (isNew = true, allOrNone) => {
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
      const role = isNew ? { operation_all: true } : { ...this.tlsOptions, ...operations };
      this.form = new KmipRoleForm(role, { isNew });
    };

    this.onSave = sinon.spy();
    this.onCancel = sinon.spy();

    // get all keys that are rendered in the operation groups
    this.operationKeys = Object.values(operationGroups()).flat();

    this.apiStub = sinon.stub(this.owner.lookup('service:api').secrets, 'kmipWriteRole').resolves();
    this.flashStub = sinon.stub(this.owner.lookup('service:flashMessages'), 'success');
    this.routerStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    this.renderComponent = (isNew, allOrNone) => {
      this.setForm(isNew, allOrNone);
      // set the roleName for existing forms
      this.roleName = isNew ? '' : 'role-1';
      return render(
        hbs`<RoleForm @roleName={{this.roleName}} @scopeName={{this.scopeName}} @form={{this.form}} @onSave={{this.onSave}} @onCancel={{this.onCancel}} />`,
        { owner: this.engine }
      );
    };
  });

  test('it should render name field and default to allow all for new roles', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.inputByAttr('name')).exists('Name field is rendered for new role');
    assert.dom(GENERAL.inputByAttr('operation_none')).isChecked('Operations are allowed by default');
    assert.dom(GENERAL.inputByAttr('operation_all')).isChecked('Allow all operations is checked by default');
  });

  test('it should hide operations when none is toggled', async function (assert) {
    await this.renderComponent();

    assert
      .dom('[data-test-kmip-section="Allowed Operations"]')
      .exists('Allowed Operations section is rendered');
    await click(GENERAL.inputByAttr('operation_none'));
    assert
      .dom('[data-test-kmip-section="Allowed Operations"]')
      .doesNotExist('Allowed Operations section is hidden');
  });

  test('it should check all operations when allow all is selected', async function (assert) {
    await this.renderComponent();

    this.operationKeys.forEach((key) => {
      assert.dom(GENERAL.inputByAttr(key)).isChecked(`${key} is checked when allow all is selected`);
    });
  });

  test('it should trigger onCancel callback', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.cancelButton);
    assert.true(this.onCancel.calledOnce, 'onCancel callback is triggered on cancel');
  });

  module('Editing existing role', function () {
    test('it should hide name field when editing', async function (assert) {
      await this.renderComponent(false);

      assert.dom(GENERAL.inputByAttr('name')).doesNotExist('Name field is hidden when editing');
    });

    test('it should populate fields when editing: operation_none', async function (assert) {
      await this.renderComponent(false, 'none');

      assert
        .dom(GENERAL.inputByAttr('operation_none'))
        .isNotChecked('Renders correct toggle state for operation_none');
      assert
        .dom('[data-test-kmip-section="Allowed Operations"]')
        .doesNotExist('Allowed Operations section is hidden');
    });

    test('it should populate fields when editing: operation_all', async function (assert) {
      await this.renderComponent(false, 'all');

      assert.dom(GENERAL.inputByAttr('operation_none')).isChecked('Allow operations toggle is checked');
      assert.dom(GENERAL.inputByAttr('operation_all')).isChecked('Allow all operations is checked');
      this.operationKeys.forEach((key) => {
        assert.dom(GENERAL.inputByAttr(key)).isChecked(`${key} is checked when allow all is selected`);
      });
    });

    test('it should populate fields when editing: selected operations', async function (assert) {
      await this.renderComponent(false);

      this.operationKeys.forEach((key) => {
        const domMethod = this.form.data[key] === true ? 'isChecked' : 'isNotChecked';
        assert.dom(GENERAL.inputByAttr(key))[domMethod](`${key} ${domMethod} correctly`);
      });
    });

    test('it should populate tls fields', async function (assert) {
      // this is a bit wonky but when the value is set by the component it's a string
      // but in order to populate it , it needs to be a number
      this.tlsOptions.tls_client_ttl = 86400;
      await this.renderComponent(false);

      assert
        .dom(GENERAL.inputByAttr('tls_client_key_bits'))
        .hasValue('521', 'TLS field is populated correctly');
      assert
        .dom(GENERAL.inputByAttr('tls_client_key_type'))
        .hasValue('ec', 'TLS field is populated correctly');
      assert
        .dom(GENERAL.ttl.input('TLS Client TTL'))
        .hasValue('1', 'TLS TTL value field is populated correctly');
      assert.dom(GENERAL.selectByAttr('ttl-unit')).hasValue('d', 'TLS TTL unit field is populated correctly');
    });

    test('it should save edited role', async function (assert) {
      await this.renderComponent(false);
      await click(GENERAL.inputByAttr('operation_none'));
      await click(GENERAL.submitButton);

      const payload = { ...this.tlsOptions, operation_none: true };
      assert.true(this.apiStub.calledWith(this.roleName, this.scopeName, this.backend, payload));
      assert.true(this.flashStub.calledWith(`Successfully saved role ${this.roleName}`));
      assert.true(this.onSave.calledOnce, 'onSave callback is triggered on save');
    });
  });

  module('Create new role', function () {
    test('it should show validation error when name is not provided', async function (assert) {
      await this.renderComponent();
      await click(GENERAL.submitButton);

      assert
        .dom(GENERAL.validationErrorByAttr('name'))
        .hasText('Name is required', 'Shows validation error for name field');
      assert.dom(GENERAL.inlineError).hasText('There is an error with this form.');
    });

    test('it should create new role without operations', async function (assert) {
      await this.renderComponent();
      await fillIn(GENERAL.inputByAttr('name'), 'new-role');
      await click(GENERAL.inputByAttr('operation_none'));
      await click(GENERAL.submitButton);

      const payload = { ...this.tlsOptions, operation_none: true };
      assert.true(
        this.apiStub.calledWith('new-role', this.scopeName, this.backend, payload),
        'Role created with no operations'
      );
      assert.true(this.flashStub.calledWith('Successfully saved role new-role'));
      assert.true(this.onSave.calledOnce, 'onSave callback is triggered on save');
    });

    test('it should create new role with all operations', async function (assert) {
      await this.renderComponent();
      await fillIn(GENERAL.inputByAttr('name'), 'new-role');
      await click(GENERAL.submitButton);

      const payload = { ...this.tlsOptions, operation_all: true };
      assert.true(
        this.apiStub.calledWith('new-role', this.scopeName, this.backend, payload),
        'Role created with all operations'
      );
      assert.true(this.flashStub.calledWith('Successfully saved role new-role'));
      assert.true(this.onSave.calledOnce, 'onSave callback is triggered on save');
    });

    test('it should create new role with selected operations', async function (assert) {
      await this.renderComponent();
      await fillIn(GENERAL.inputByAttr('name'), 'new-role');
      await click(GENERAL.inputByAttr('operation_all'));
      await click(GENERAL.inputByAttr('operation_decrypt'));
      await click(GENERAL.inputByAttr('operation_create'));
      await click(GENERAL.inputByAttr('operation_get_attributes'));
      await click(GENERAL.submitButton);

      const payload = {
        ...this.tlsOptions,
        operation_decrypt: true,
        operation_create: true,
        operation_get_attributes: true,
      };
      assert.true(
        this.apiStub.calledWith('new-role', this.scopeName, this.backend, payload),
        'Role created with all operations'
      );
      assert.true(this.flashStub.calledWith('Successfully saved role new-role'));
      assert.true(this.onSave.calledOnce, 'onSave callback is triggered on save');
    });
  });
});
