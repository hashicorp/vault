/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import LdapStaticRoleForm from 'vault/forms/secrets/ldap/roles/static';
import LdapDynamicRoleForm from 'vault/forms/secrets/ldap/roles/dynamic';
import { formatError, overrideResponse } from 'vault/tests/helpers/stubs';

module('Integration | Component | ldap | Page::Role::CreateAndEdit', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.backend = 'ldap-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    const routerStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    this.transitionCalledWith = (routeName, name) => {
      const route = `vault.cluster.secrets.backend.ldap.${routeName}`;
      const args = name ? [route, name] : [route];
      return routerStub.calledWith(...args);
    };

    const [staticRoleData, dynamicRoleData] = ['static', 'dynamic'].map((roleType) => {
      const name = `${roleType}-role`;
      const role = this.server.create('ldap-role', roleType, { name });
      delete role.id;
      delete role.type;
      return role;
    });

    this.createModel = {
      staticForm: new LdapStaticRoleForm({}, { isNew: true }),
      dynamicForm: new LdapDynamicRoleForm({ default_ttl: '1h', max_ttl: '24h' }, { isNew: true }),
    };
    this.staticEditModel = {
      staticForm: new LdapStaticRoleForm(staticRoleData),
    };
    this.dynamicEditModel = {
      dynamicForm: new LdapDynamicRoleForm(dynamicRoleData),
    };

    this.breadcrumbs = [
      { label: 'ldap', route: 'overview' },
      { label: 'Roles', route: 'roles' },
      { label: 'Create' },
    ];

    this.renderComponent = () => {
      return render(
        hbs`<Page::Role::CreateAndEdit @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`,
        { owner: this.engine }
      );
    };
  });

  test('it should display different form fields based on type', async function (assert) {
    assert.expect(12);

    this.model = this.createModel;
    await this.renderComponent();

    assert.dom('[data-test-radio-card="static"]').isChecked('Static role type selected by default');

    const checkFields = (fields) => {
      fields.forEach((field) => {
        assert.dom(GENERAL.fieldByAttr(field)).exists(`${field} field renders when static type is selected`);
      });
    };

    checkFields(['name', 'dn', 'username', 'rotation_period']);
    await click('[data-test-radio-card="dynamic"]');
    checkFields([
      'name',
      'default_ttl',
      'max_ttl',
      'username_template',
      'creation_ldif',
      'deletion_ldif',
      'rollback_ldif',
    ]);
  });

  test('it should populate form and disable type cards when editing', async function (assert) {
    assert.expect(17);

    const checkFields = (fields, element = 'input:last-child') => {
      fields.forEach((field) => {
        const isLdif = field.includes('ldif');
        const method = isLdif ? 'includesText' : 'hasValue';
        const value = isLdif ? 'dn: cn={{.Username}},ou=users,dc=learn,dc=example' : this.form.data[field];
        assert.dom(`${GENERAL.fieldByAttr(field)} ${element}`)[method](value, `${field} field value renders`);
      });
    };
    const checkTtl = (fields) => {
      fields.forEach((field) => {
        assert
          .dom(`${GENERAL.fieldByAttr(field)} [data-test-ttl-inputs] input`)
          .hasAnyValue(`${field} field ttl value renders`);
      });
    };

    this.model = this.staticEditModel;
    this.form = this.model.staticForm;
    await this.renderComponent();
    assert.dom('[data-test-radio-card="static"]').isChecked('Type is set when editing role');
    assert.dom('[data-test-radio-card="static"]').isDisabled('Type selection is disabled when editing');
    assert.dom(GENERAL.inputByAttr('name')).isDisabled('Name field is disabled when editing');
    checkFields(['name', 'dn', 'username']);
    checkTtl(['rotation_period']);

    this.model = this.dynamicEditModel;
    this.form = this.model.dynamicForm;
    await this.renderComponent();
    assert.dom('[data-test-radio-card="dynamic"]').isChecked('Type is set when editing role');
    assert.dom('[data-test-radio-card="dynamic"]').isDisabled('Type selection is disabled when editing');
    assert.dom(GENERAL.inputByAttr('name')).isDisabled('Name field is disabled when editing');
    checkFields(['name', 'username_template']);
    checkTtl(['default_ttl', 'max_ttl']);
    checkFields(['creation_ldif', 'deletion_ldif', 'rollback_ldif'], '.cm-content');
  });

  test('it should go back to list route on cancel', async function (assert) {
    this.model = this.staticEditModel;

    await this.renderComponent();
    await click(GENERAL.cancelButton);

    assert.ok(this.transitionCalledWith('roles'), 'Transitions to roles list route on cancel');
  });

  test('it should validate form fields', async function (assert) {
    const submitAndAssert = async (fields) => {
      await click(GENERAL.submitButton);
      fields.forEach((field) => {
        assert.dom(GENERAL.validationErrorByAttr(field)).exists('Validation message renders');
      });

      assert
        .dom('[data-test-invalid-form-message]')
        .hasText(`There are ${fields.length} errors with this form.`);
    };

    this.model = this.createModel;
    await this.renderComponent();
    await submitAndAssert(['name', 'username', 'rotation_period']);

    await click('[data-test-radio-card="dynamic"]');
    await submitAndAssert(['name', 'creation_ldif', 'deletion_ldif']);
  });

  test('it should create new static role', async function (assert) {
    assert.expect(2);

    const writeStub = sinon.stub(this.owner.lookup('service:api').secrets, 'ldapWriteStaticRole').resolves();

    this.model = this.createModel;
    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('name'), 'test-role');
    await fillIn(GENERAL.inputByAttr('dn'), 'foo');
    await fillIn(GENERAL.inputByAttr('username'), 'bar');
    await fillIn(GENERAL.ttl.input('Rotation period'), 5);
    await click(GENERAL.submitButton);

    const payload = { dn: 'foo', username: 'bar', rotation_period: '5s' };
    assert.true(
      writeStub.calledWith('test-role', this.backend, payload),
      'Request made to create static role with correct properties'
    );
    assert.true(
      this.transitionCalledWith('roles.role.details', 'static', 'test-role'),
      'Transitions to role details route on save success'
    );
  });

  test('it should save edited role with correct properties', async function (assert) {
    assert.expect(2);

    const writeStub = sinon.stub(this.owner.lookup('service:api').secrets, 'ldapWriteStaticRole').resolves();
    this.server.post('/ldap-test/static-role/:name', (schema, req) => {
      const data = JSON.parse(req.requestBody);
      const expected = { dn: 'foo', username: 'bar', rotation_period: '30s' };
      assert.deepEqual(expected, data, 'POST request made to save role with correct properties');
    });

    this.model = this.staticEditModel;
    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('dn'), 'foo');
    await fillIn(GENERAL.inputByAttr('username'), 'bar');
    await fillIn(GENERAL.ttl.input('Rotation period'), 30);
    await click(GENERAL.submitButton);

    const payload = { dn: 'foo', username: 'bar', rotation_period: '30s' };
    assert.true(
      writeStub.calledWith('static-role', this.backend, payload),
      'Request made to edit role with correct properties'
    );
    assert.true(
      this.transitionCalledWith('roles.role.details', 'static', 'test-role'),
      'Transitions to role details route on save success'
    );
  });

  test('it should make a request to correct endpoint for dynamic roles', async function (assert) {
    assert.expect(2);

    const writeStub = sinon.stub(this.owner.lookup('service:api').secrets, 'ldapWriteDynamicRole').resolves();

    this.model = this.dynamicEditModel;
    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('username_template'), 'bar');
    await click(GENERAL.submitButton);

    const { name, ...data } = this.model.dynamicForm.data;
    const payload = { ...data, username_template: 'bar' };
    assert.true(
      writeStub.calledWith(name, this.backend, payload),
      'Request made to correct endpoint for dynamic role'
    );
    assert.true(
      this.transitionCalledWith('roles.role.details', 'dynamic', name),
      'Transitions to role details route on save success'
    );
  });

  test('it should display api error when creating static roles fails', async function (assert) {
    this.server.post('/:backend/static-role/:name', () => {
      return overrideResponse(500, formatError('uh oh!'));
    });
    this.model = this.createModel;
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('name'), 'test-role');
    await fillIn(GENERAL.inputByAttr('dn'), 'foo');
    await fillIn(GENERAL.inputByAttr('username'), 'bar');
    await fillIn(GENERAL.ttl.input('Rotation period'), 5);
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.messageError)
      .hasText('Error uh oh!', 'it renders error message returned from the API');
  });

  test('it should display api error when creating dynamic roles fails', async function (assert) {
    this.server.post('/:backend/role/:name', () => {
      return overrideResponse(500, formatError('uh oh!'));
    });
    this.model = this.dynamicEditModel;
    await this.renderComponent();
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.messageError)
      .hasText('Error uh oh!', 'it renders error message returned from the API');
  });
});
