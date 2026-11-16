/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, fillIn } from '@ember/test-helpers';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import DatabaseRoleForm, { ROLE_DEFAULTS, normalizeDatabaseRole } from 'vault/forms/database/role';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Integration | Component | database-role-edit', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.version = this.owner.lookup('service:version');
    this.router = this.owner.lookup('service:router');
    this.router.transitionTo = sinon.stub();
    const staticRole = {
      database: ['my-mongodb-database'],
      backend: 'database',
      username: 'staticTestUser',
      skip_import_rotation: false,
      type: 'static',
      name: 'my-static-role',
    };
    const dynamicRole = {
      database: ['my-mongodb-database'],
      backend: 'database',
      type: 'dynamic',
      name: 'my-dynamic-role',
    };
    // the edit view changes the role so it gets a form, the show view is read-only plain data
    this.editForm = (role) => new DatabaseRoleForm({ ...role }, { isNew: false });
    this.createForm = (role) => new DatabaseRoleForm({ ...ROLE_DEFAULTS, ...role }, { isNew: true });
    this.showData = (role) =>
      normalizeDatabaseRole({ ...role, db_name: role.database[0] }, role.name, 'database', role.type);
    this.staticRole = staticRole;
    this.dynamicRole = dynamicRole;
    this.root = {
      label: 'database',
      text: 'database',
      path: 'vault.cluster.secrets.backend.list-root',
      model: 'database',
    };
    // saving a new role also adds it to the allowed roles of its connection
    this.server.get('/database/config/my-mongodb-database', () => ({
      data: { plugin_name: 'mongodb-database-plugin', allowed_roles: [] },
    }));
    this.server.post('/database/config/my-mongodb-database', () => new Response(null, { status: 204 }));
  });

  test('it should display form errors when trying to create a role without required fields', async function (assert) {
    this.form = this.createForm({ backend: 'database', database: ['my-mongodb-database'], type: 'static' });

    await render(hbs`<DatabaseRoleEdit @model={{this.form}} @root={{this.root}} @mode="create"/>`);
    await click(GENERAL.submitButton);

    assert.dom('[data-test-inline-error-message]').exists('Inline form errors exist');
  });

  test('it should let user edit a static role and only sends the changed attributes with the username', async function (assert) {
    assert.expect(3);
    // a role read from the API always has a rotation period
    this.form = this.editForm({ ...this.staticRole, rotation_period: 86400 });

    this.server.post(`/database/static-roles/my-static-role`, (schema, req) => {
      assert.true(true, 'request made to update static role');
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          username: 'staticTestUser',
          rotation_period: '1728000s', // 20 days in seconds
        },
        'it updates static role with correct payload and never sends a password'
      );
    });

    await render(hbs`<DatabaseRoleEdit @model={{this.form}} @root={{this.root}} @mode="edit"/>`);
    await fillIn('[data-test-ttl-value="Rotation period"]', '20');
    await click(GENERAL.submitButton);
    assert.true(
      this.router.transitionTo.calledWith('vault.cluster.secrets.backend.show', 'role/my-static-role')
    );
  });

  test('it should create a dynamic role with the connection as db_name and the creation statement as a list', async function (assert) {
    assert.expect(2);
    this.form = this.createForm({ ...this.dynamicRole, name: 'new-dynamic' });
    this.form.creation_statement = '{"db":"admin","roles":[{"role":"readWrite"}]}';

    this.server.post(`/database/roles/new-dynamic`, (schema, req) => {
      assert.true(true, 'request made to create dynamic role');
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          db_name: 'my-mongodb-database',
          default_ttl: '1h',
          max_ttl: '24h',
          creation_statements: ['{"db":"admin","roles":[{"role":"readWrite"}]}'],
        },
        'it maps the connection to db_name and the single statement to a list'
      );
    });

    await render(hbs`<DatabaseRoleEdit @model={{this.form}} @root={{this.root}} @mode="create"/>`);
    await click(GENERAL.submitButton);
  });

  test('it should stop saving a new role when updating the connection requires control group approval', async function (assert) {
    this.form = this.createForm({ ...this.dynamicRole, name: 'new-dynamic' });
    const controlGroup = this.owner.lookup('service:control-group');
    const saveToken = sinon.stub(controlGroup, 'saveTokenFromError');
    sinon.stub(controlGroup, 'logFromError').returns({ content: 'A Control Group was encountered' });
    sinon
      .stub(this.owner.lookup('service:api').request, 'post')
      .rejects(getErrorResponse({ accessor: 'foobar', isControlGroupError: true }, 403));
    this.server.post('/database/roles/new-dynamic', () => {
      assert.true(false, 'role is not saved before the connection update is authorized');
    });

    await render(hbs`<DatabaseRoleEdit @model={{this.form}} @root={{this.root}} @mode="create"/>`);
    await click(GENERAL.submitButton);

    assert.true(saveToken.calledOnce, 'control group token is saved');
    assert.dom(GENERAL.messageError).includesText('A Control Group was encountered');
    assert.false(this.router.transitionTo.called, 'it does not transition');
  });

  test('enterprise: it should successfully create user that does not rotate immediately', async function (assert) {
    assert.expect(4);
    this.version.type = 'enterprise';
    this.form = this.createForm(this.staticRole);

    this.server.post(`/database/static-roles/my-static-role`, (schema, req) => {
      assert.true(true, 'request made to create static role');
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          db_name: 'my-mongodb-database',
          username: 'staticTestUser',
          password: 'testPassword',
          rotation_period: '172800s', // 2 days in seconds
          skip_import_rotation: true,
        },
        'it creates a static role with correct payload'
      );
    });

    await render(hbs`<DatabaseRoleEdit @model={{this.form}} @root={{this.root}} @mode="create"/>`);
    await fillIn(GENERAL.ttl.input('Rotation period'), '2');
    await click(GENERAL.toggleInput('toggle-skip_import_rotation'));
    await fillIn(GENERAL.inputByAttr('password'), 'testPassword'); // fill in password field
    await click(GENERAL.submitButton);

    this.role = this.showData({ ...this.staticRole, skip_import_rotation: true });
    await render(hbs`<DatabaseRoleEdit @model={{this.role}} @root={{this.root}} @mode="show"/>`);
    assert.dom(GENERAL.infoRowValue('Rotate immediately')).containsText('No');
    assert.dom(GENERAL.infoRowValue('password')).doesNotExist(); // verify password field doesn't show on details view
  });

  test('enterprise: a new static role inherits skip_static_role_import_rotation from its connection', async function (assert) {
    assert.expect(2);
    this.version.type = 'enterprise';
    this.server.get('/database/config/my-mongodb-database', () => ({
      data: {
        plugin_name: 'mongodb-database-plugin',
        allowed_roles: [],
        skip_static_role_import_rotation: true,
      },
    }));
    const { skip_import_rotation, ...role } = this.staticRole; // eslint-disable-line no-unused-vars
    this.form = this.createForm(role);

    this.server.post(`/database/static-roles/my-static-role`, (schema, req) => {
      assert.true(JSON.parse(req.requestBody).skip_import_rotation, 'payload keeps the connection setting');
    });

    await render(hbs`<DatabaseRoleEdit @model={{this.form}} @root={{this.root}} @mode="create"/>`);
    assert.dom(GENERAL.toggleInput('toggle-skip_import_rotation')).isNotChecked('rotate immediately is off');
    await fillIn(GENERAL.inputByAttr('password'), 'testPassword');
    await click(GENERAL.submitButton);
  });

  test('enterprise: it should show edit button for password when role has not been rotated', async function (assert) {
    this.version.type = 'enterprise';
    this.form = this.editForm(this.staticRole);

    await render(hbs`<DatabaseRoleEdit @model={{this.form}} @root={{this.root}} @mode="edit"/>`);
    assert.dom(GENERAL.icon('edit')).exists(); // verify password field is enabled for edit & enable button is rendered bc role hasn't been rotated
  });

  test('enterprise: it should not show edit button for password when role has been rotated', async function (assert) {
    this.version.type = 'enterprise';
    // Setting a sample rotation time here to simulate what returns from BE after rotation
    this.form = this.editForm({
      ...this.staticRole,
      last_vault_rotation: '2025-04-21T12:51:59.063124-04:00',
    });

    await render(hbs`<DatabaseRoleEdit @model={{this.form}} @root={{this.root}} @mode="edit"/>`);
    assert.dom(GENERAL.icon('edit')).doesNotExist(); // verify password field is disabled for edit & enable button isn't rendered bc role has already been rotated
  });

  test('enterprise: it should successfully create user that does rotate immediately & verify warning modal pops up', async function (assert) {
    this.version.type = 'enterprise';
    this.form = this.createForm(this.staticRole);

    await render(hbs`<DatabaseRoleEdit @model={{this.form}} @root={{this.root}} @mode="create"/>`);
    await click(GENERAL.submitButton);

    assert.dom('[data-test-issuer-warning]').exists(); // check if warning modal shows after clicking save
    await click('[data-test-issuer-save]'); // click continue button on modal

    this.role = this.showData(this.staticRole);
    await render(hbs`<DatabaseRoleEdit @model={{this.role}} @root={{this.root}} @mode="show"/>`);
    assert.dom(GENERAL.infoRowValue('Rotate immediately')).containsText('Yes');
  });

  test('it should show Get credentials button when a user has the correct policy', async function (assert) {
    this.role = this.showData(this.staticRole);
    this.capabilities = { canGetCredentials: true };
    await render(
      hbs`<DatabaseRoleEdit @model={{this.role}} @capabilities={{this.capabilities}} @root={{this.root}} @mode="show"/>`
    );
    assert.dom(GENERAL.button('static')).exists('Get credentials button exists');
  });

  test('it should show Generate credentials button when a user has the correct policy', async function (assert) {
    this.role = this.showData(this.dynamicRole);
    this.capabilities = { canGenerateCredentials: true };
    await render(
      hbs`<DatabaseRoleEdit @model={{this.role}} @capabilities={{this.capabilities}} @root={{this.root}} @mode="show"/>`
    );
    assert.dom(GENERAL.button('dynamic')).exists('Generate credentials button exists');
  });

  test('it should not show credential or management actions without capabilities', async function (assert) {
    this.role = this.showData(this.staticRole);
    this.capabilities = {};
    await render(
      hbs`<DatabaseRoleEdit @model={{this.role}} @capabilities={{this.capabilities}} @root={{this.root}} @mode="show"/>`
    );
    assert.dom(GENERAL.button('static')).doesNotExist('Get credentials button is hidden');
    assert.dom('[data-test-database-role-delete]').doesNotExist('Delete is hidden');
    assert.dom('[data-test-edit-link]').doesNotExist('Edit is hidden');
  });
});
