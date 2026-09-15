/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, fillIn, findAll } from '@ember/test-helpers';
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
    this.transitionCalledWith = (routeName, ...extraArgs) => {
      const route = `vault.cluster.secrets.backend.ldap.${routeName}`;
      return routerStub.calledWith(route, ...extraArgs);
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
      isSelfManaged: false,
    };
    this.selfManagedCreateModel = {
      staticForm: new LdapStaticRoleForm({}, { isNew: true, isSelfManaged: true }),
      dynamicForm: new LdapDynamicRoleForm({ default_ttl: '1h', max_ttl: '24h' }, { isNew: true }),
      isSelfManaged: true,
    };
    this.staticEditModel = {
      staticForm: new LdapStaticRoleForm(staticRoleData),
      isSelfManaged: false,
    };
    this.selfManagedEditModel = {
      staticForm: new LdapStaticRoleForm(staticRoleData, { isSelfManaged: true }),
      isSelfManaged: true,
    };
    this.dynamicEditModel = {
      dynamicForm: new LdapDynamicRoleForm(dynamicRoleData),
      isSelfManaged: false,
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

      assert.dom(GENERAL.invalidFormMessage).hasText(`There are ${fields.length} errors with this form.`);
    };

    this.model = this.createModel;
    await this.renderComponent();
    await submitAndAssert(['name', 'username']);

    await click('[data-test-radio-card="dynamic"]');
    await submitAndAssert(['name', 'creation_ldif', 'deletion_ldif']);
  });

  test('it should render the LDAP validation messages for an empty static role form', async function (assert) {
    this.model = this.createModel;
    await this.renderComponent();

    await click(GENERAL.submitButton);

    assert
      .dom(GENERAL.validationErrorByAttr('name'))
      .hasText('Role name is required', 'name renders the LDAP-specific message');
    assert
      .dom(GENERAL.validationErrorByAttr('username'))
      .hasText('Enter the username of the LDAP account', 'username renders the LDAP-specific message');
    assert
      .dom(GENERAL.validationErrorByAttr('rotation_period'))
      .doesNotExist('rotation period never reports a validation error');
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

    this.model = this.staticEditModel;
    await this.renderComponent();

    // dn is editDisabled, so it cannot be changed here; it keeps its seeded value
    const seededDn = this.staticEditModel.staticForm.data.dn;
    await fillIn(GENERAL.inputByAttr('username'), 'bar');
    await fillIn(GENERAL.ttl.input('Rotation period'), 30);
    await click(GENERAL.submitButton);

    const payload = { dn: seededDn, username: 'bar', rotation_period: '30s' };
    assert.true(
      writeStub.calledWith('static-role', this.backend, payload),
      'Request made to edit role with correct properties'
    );
    assert.true(
      this.transitionCalledWith('roles.role.details', 'static', 'static-role'),
      'Transitions to role details route on save success'
    );
  });

  test('it should disable the distinguished name field only when editing', async function (assert) {
    this.model = this.createModel;
    await this.renderComponent();
    assert.dom(GENERAL.inputByAttr('dn')).isNotDisabled('dn is editable when creating a role');

    this.model = this.staticEditModel;
    await this.renderComponent();
    assert.dom(GENERAL.inputByAttr('dn')).isDisabled('dn is read-only when editing an existing role');
  });

  test('it should require dn and password only for self-managed static roles', async function (assert) {
    const rootManaged = new LdapStaticRoleForm({}, { isNew: true });
    const selfManaged = new LdapStaticRoleForm({}, { isNew: true, isSelfManaged: true });
    const optionsFor = (form, name) => form.formFields.find((field) => field.name === name).options;

    assert.notOk(optionsFor(rootManaged, 'dn').isRequired, 'dn is not required when root-managed');
    assert.notOk(
      optionsFor(rootManaged, 'password').isRequired,
      'password is not required when root-managed'
    );
    assert.true(optionsFor(selfManaged, 'dn').isRequired, 'dn is required when self-managed');
    assert.true(optionsFor(selfManaged, 'password').isRequired, 'password is required when self-managed');

    // name and username are required in both modes, matching the backend contract
    ['name', 'username'].forEach((field) => {
      assert.true(optionsFor(rootManaged, field).isRequired, `${field} is required when root-managed`);
      assert.true(optionsFor(selfManaged, field).isRequired, `${field} is required when self-managed`);
    });
  });

  test('it should switch the username subText based on self-managed', async function (assert) {
    const rootManaged = new LdapStaticRoleForm({}, { isNew: true });
    const selfManaged = new LdapStaticRoleForm({}, { isNew: true, isSelfManaged: true });
    const subTextFor = (form) => form.formFields.find((field) => field.name === 'username').options.subText;

    assert.strictEqual(
      subTextFor(rootManaged),
      "The name of the user to be used when logging in. This is useful when DN isn't used for login purposes.",
      'root-managed username subText'
    );
    assert.strictEqual(
      subTextFor(selfManaged),
      'The username of the existing LDAP entry to manage.',
      'self-managed username subText'
    );
  });

  test('it should validate dn and password only for self-managed static roles', async function (assert) {
    const rootManaged = new LdapStaticRoleForm({ name: 'test-role', username: 'bar' }, { isNew: true });
    assert.true(rootManaged.toJSON().isValid, 'form is valid without dn or password when root-managed');

    const selfManaged = new LdapStaticRoleForm(
      { name: 'test-role', username: 'bar' },
      { isNew: true, isSelfManaged: true }
    );
    const { isValid, state, data } = selfManaged.toJSON();
    assert.false(isValid, 'form is invalid without dn or password when self-managed');
    assert.deepEqual(
      state.dn.errors,
      ['Enter the distinguished name (DN) of the LDAP account'],
      'dn renders the LDAP-specific message'
    );
    assert.deepEqual(
      state.password.errors,
      ['Enter the current password for this LDAP account'],
      'password renders the LDAP-specific message'
    );
    assert.deepEqual(
      Object.keys(data).sort(),
      ['name', 'username'],
      'validation adds nothing to the payload'
    );

    const complete = new LdapStaticRoleForm(
      { name: 'test-role', username: 'bar', dn: 'cn=test', password: 'secret' },
      { isNew: true, isSelfManaged: true }
    );
    assert.true(complete.toJSON().isValid, 'form is valid once dn and password are supplied');
  });

  test('it should not require a password when editing an existing self-managed role', async function (assert) {
    const writeStub = sinon.stub(this.owner.lookup('service:api').secrets, 'ldapWriteStaticRole').resolves();

    this.model = this.selfManagedEditModel;
    await this.renderComponent();

    assert
      .dom(GENERAL.enableField('password'))
      .exists('password stays locked behind Enable input, so it is not presented as required');

    await click(GENERAL.submitButton);

    assert
      .dom(GENERAL.invalidFormMessage)
      .doesNotExist('leaving the stored password untouched does not block the save');
    const [, , payload] = writeStub.getCall(0).args;
    assert.false(
      Object.prototype.hasOwnProperty.call(payload, 'password'),
      'password is absent from the payload, so the stored value is left untouched'
    );
  });

  test('it should clear the rotation period to zero rather than flagging it as required', async function (assert) {
    const writeStub = sinon.stub(this.owner.lookup('service:api').secrets, 'ldapWriteStaticRole').resolves();

    this.model = this.selfManagedCreateModel;
    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('name'), 'test-role');
    await fillIn(GENERAL.inputByAttr('username'), 'bar');
    await fillIn(GENERAL.inputByAttr('dn'), 'cn=test');
    await fillIn(GENERAL.inputByAttr('password'), 'secret');
    await fillIn(GENERAL.ttl.input('Rotation period'), '');

    assert
      .dom('.ttl-value-error')
      .doesNotExist('an empty rotation period is a valid zero, not a required-field error');

    await click(GENERAL.submitButton);

    const [, , payload] = writeStub.getCall(0).args;
    assert.strictEqual(payload.rotation_period, '0s', 'the cleared field saves as zero');
  });

  test('it should render an editable password field on the create page', async function (assert) {
    this.model = this.createModel;
    await this.renderComponent();

    assert.dom(GENERAL.fieldByAttr('password')).exists('password field renders when creating a role');
    assert
      .dom(GENERAL.enableField('password'))
      .doesNotExist('password is directly editable on create, not locked behind an Enable input button');
    assert
      .dom(GENERAL.inputByAttr('password'))
      .isNotDisabled('password is editable')
      .hasAttribute('type', 'password', 'password input masks what the user types');
  });

  // On edit the stored password is never sent to the browser. EnableInput renders a disabled,
  // masked placeholder until the user explicitly opts in to replacing the value.
  test('it should lock the password field behind an Enable input button on the edit page', async function (assert) {
    this.model = this.staticEditModel;
    await this.renderComponent();

    assert.dom(GENERAL.enableField('password')).exists('Enable input button renders');
    assert
      .dom(GENERAL.inputByAttr('password'))
      .isDisabled('the placeholder cannot be edited until Enable input is clicked')
      .hasValue('**********', 'the stored password is masked, never revealed');
  });

  test('it should reveal an empty editable password field when Enable input is clicked', async function (assert) {
    this.model = this.staticEditModel;
    await this.renderComponent();

    await click(GENERAL.enableField('password'));

    assert.dom(GENERAL.enableField('password')).doesNotExist('Enable input button is replaced by the field');
    assert
      .dom(GENERAL.inputByAttr('password'))
      .isNotDisabled('password is now editable')
      .hasAttribute('type', 'password', 'password input masks what the user types')
      .hasNoValue('the field starts empty so the stored password is never exposed');
  });

  test('it should omit password from the payload when Enable input was never clicked', async function (assert) {
    const writeStub = sinon.stub(this.owner.lookup('service:api').secrets, 'ldapWriteStaticRole').resolves();

    this.model = this.staticEditModel;
    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('username'), 'bar');
    await click(GENERAL.submitButton);

    const [, , payload] = writeStub.getCall(0).args;
    assert.false(
      Object.prototype.hasOwnProperty.call(payload, 'password'),
      'password is absent from the payload, so the stored value is left untouched'
    );
  });

  test('it should include password in the payload after Enable input and a new value', async function (assert) {
    const writeStub = sinon.stub(this.owner.lookup('service:api').secrets, 'ldapWriteStaticRole').resolves();

    this.model = this.staticEditModel;
    await this.renderComponent();

    await click(GENERAL.enableField('password'));
    await fillIn(GENERAL.inputByAttr('password'), 'new-secret');
    await fillIn(GENERAL.inputByAttr('username'), 'bar');
    await click(GENERAL.submitButton);

    const [, , payload] = writeStub.getCall(0).args;
    assert.strictEqual(payload.password, 'new-secret', 'the new password is sent');
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

  // ––––– self-managed mounts –––––

  test('it should hide the role type cards when the mount is self-managed', async function (assert) {
    this.model = this.createModel;
    await this.renderComponent();
    assert.dom(GENERAL.radioCardByAttr('static')).exists('static card renders for a root-managed mount');
    assert.dom(GENERAL.radioCardByAttr('dynamic')).exists('dynamic card renders for a root-managed mount');

    this.model = this.selfManagedCreateModel;
    await this.renderComponent();
    assert
      .dom(GENERAL.radioCardByAttr('static'))
      .doesNotExist('static card is hidden for a self-managed mount');
    assert
      .dom(GENERAL.radioCardByAttr('dynamic'))
      .doesNotExist('dynamic role type is not offered for a self-managed mount');
  });

  test('it should hide the role type cards when editing a self-managed role', async function (assert) {
    this.model = this.selfManagedEditModel;
    await this.renderComponent();

    assert.dom(GENERAL.radioCardByAttr('static')).doesNotExist('static card is hidden on edit');
    assert.dom(GENERAL.radioCardByAttr('dynamic')).doesNotExist('dynamic card is hidden on edit');
  });

  test('it should save a self-managed static role', async function (assert) {
    const writeStub = sinon.stub(this.owner.lookup('service:api').secrets, 'ldapWriteStaticRole').resolves();

    this.model = this.selfManagedCreateModel;
    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('name'), 'test-role');
    await fillIn(GENERAL.inputByAttr('dn'), 'cn=test,dc=example,dc=com');
    await fillIn(GENERAL.inputByAttr('username'), 'bar');
    await fillIn(GENERAL.inputByAttr('password'), 'secret');
    await click(GENERAL.submitButton);

    assert.true(writeStub.calledOnce, 'the static role endpoint is called');
    const [name, , payload] = writeStub.getCall(0).args;
    assert.strictEqual(name, 'test-role', 'role name is sent');
    assert.strictEqual(payload.dn, 'cn=test,dc=example,dc=com', 'dn is sent');
    assert.strictEqual(payload.password, 'secret', 'password is sent');
  });

  // dn and password are only enforced for self-managed mounts, where Vault manages an existing
  // LDAP account rather than one it created itself.
  test('it should require dn and password only when the mount is self-managed', async function (assert) {
    this.model = this.selfManagedCreateModel;
    await this.renderComponent();
    await click(GENERAL.submitButton);

    assert
      .dom(GENERAL.validationErrorByAttr('dn'))
      .hasText('Enter the distinguished name (DN) of the LDAP account', 'dn is required');
    assert
      .dom(GENERAL.validationErrorByAttr('password'))
      .hasText('Enter the current password for this LDAP account', 'password is required');

    this.model = this.createModel;
    await this.renderComponent();
    await click(GENERAL.submitButton);

    assert.dom(GENERAL.validationErrorByAttr('dn')).doesNotExist('dn is optional for a root-managed mount');
    assert
      .dom(GENERAL.validationErrorByAttr('password'))
      .doesNotExist('password is optional for a root-managed mount');
  });

  test('it should switch the username helper text based on the mount type', async function (assert) {
    this.model = this.createModel;
    await this.renderComponent();
    assert
      .dom(`${GENERAL.fieldByAttr('username')} ${GENERAL.helpText}`)
      .hasText(
        "The name of the user to be used when logging in. This is useful when DN isn't used for login purposes.",
        'root-managed helper text'
      );

    this.model = this.selfManagedCreateModel;
    await this.renderComponent();
    assert
      .dom(`${GENERAL.fieldByAttr('username')} ${GENERAL.helpText}`)
      .hasText('The username of the existing LDAP entry to manage.', 'self-managed helper text');
  });

  test('it should render the static role fields in the order set by the design spec', async function (assert) {
    const renderedOrder = async (model) => {
      this.model = model;
      await this.renderComponent();
      return findAll('[data-test-field]').map((el) => el.getAttribute('data-test-field'));
    };

    const expected = ['name', 'dn', 'username', 'password', 'rotation_period'];
    assert.deepEqual(await renderedOrder(this.createModel), expected, 'root-managed field order');
    assert.deepEqual(await renderedOrder(this.selfManagedCreateModel), expected, 'self-managed field order');
  });

  test('it should title the page "Create static role" only when creating on a self-managed mount', async function (assert) {
    this.model = this.selfManagedCreateModel;
    await this.renderComponent();
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Create static role', 'self-managed create title');

    this.model = this.createModel;
    await this.renderComponent();
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Create Role', 'root-managed create title is unchanged');

    this.model = this.selfManagedEditModel;
    await this.renderComponent();
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Edit Role', 'edit title is unchanged');
  });
});
