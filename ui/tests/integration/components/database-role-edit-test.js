/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { capabilitiesStub } from 'vault/tests/helpers/stubs';
import { click, fillIn } from '@ember/test-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | database-role-edit', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.version = this.owner.lookup('service:version');
    this.store = this.owner.lookup('service:store');
    this.store.pushPayload('database-role', {
      modelName: 'database/role',
      database: ['my-mongodb-database'],
      backend: 'database',
      username: 'staticTestUser',
      skip_import_rotation: false,
      type: 'static',
      name: 'my-static-role',
      id: 'my-static-role',
    });
    this.store.pushPayload('database-role', {
      modelName: 'database/role',
      database: ['my-mongodb-database'],
      backend: 'database',
      type: 'dynamic',
      name: 'my-dynamic-role',
      id: 'my-dynamic-role',
    });
    this.store.pushPayload('database-role', {
      modelName: 'database/role',
      database: ['my-mongodb-database'],
      id: 'test-role',
      type: 'static',
      name: 'test-role',
    });
    this.modelStatic = this.store.peekRecord('database/role', 'my-static-role');
    this.modelDynamic = this.store.peekRecord('database/role', 'my-dynamic-role');
    this.modelEmpty = this.store.peekRecord('database/role', 'test-role');
  });

  test('it should display form errors when trying to create a role without required fields', async function (assert) {
    this.server.post('/sys/capabilities-self', capabilitiesStub('database/static-creds/my-role', ['create']));

    await render(hbs`<DatabaseRoleEdit @model={{this.modelEmpty}} @mode="create"/>`);
    await click(GENERAL.submitButton);

    assert.dom('[data-test-inline-error-message]').exists('Inline form errors exist');
  });

  test('it should let user edit a static role when given update capability', async function (assert) {
    this.server.post('/sys/capabilities-self', capabilitiesStub('database/static-creds/my-role', ['update']));

    this.server.post(`/database/static-roles/my-static-role`, (schema, req) => {
      assert.true(true, 'request made to update static role');
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          username: 'staticTestUser',
          rotation_period: '1728000s', // 20 days in seconds
        },
        'it updates static role with correct payload'
      );
    });

    await render(hbs`<DatabaseRoleEdit @model={{this.modelStatic}} @mode="edit"/>`);
    await fillIn('[data-test-ttl-value="Rotation period"]', '20');
    await click(GENERAL.submitButton);
  });

  test('enterprise: it should successfully create user that does not rotate immediately', async function (assert) {
    this.version.type = 'enterprise';
    this.server.post('/sys/capabilities-self', capabilitiesStub('database/static-creds/my-role', ['create']));
    this.server.post(`/database/static-roles/my-static-role`, (schema, req) => {
      assert.true(true, 'request made to create static role');
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          path: 'static-roles',
          username: 'staticTestUser',
          password: 'testPassword',
          rotation_period: '172800s', // 2 days in seconds
          skip_import_rotation: true,
        },
        'it creates a static role with correct payload'
      );
    });

    await render(hbs`<DatabaseRoleEdit @model={{this.modelStatic}} @mode="create"/>`);
    await fillIn(GENERAL.ttl.input('Rotation period'), '2');
    await click(GENERAL.toggleInput('toggle-skip_import_rotation'));
    await fillIn(GENERAL.inputByAttr('password'), 'testPassword'); // fill in password field

    await click(GENERAL.submitButton);

    await render(hbs`<DatabaseRoleEdit @model={{this.modelStatic}} @mode="show"/>`);
    assert.dom(GENERAL.infoRowValue('Rotate immediately')).containsText('No');
    assert.dom(GENERAL.infoRowValue('password')).doesNotExist(); // verify password field doesn't show on details view
  });

  test('enterprise: it should show edit button for password when role has not been rotated', async function (assert) {
    this.version.type = 'enterprise';
    this.server.post('/sys/capabilities-self', capabilitiesStub('database/static-creds/my-role', ['update']));

    await render(hbs`<DatabaseRoleEdit @model={{this.modelStatic}} @mode="edit"/>`);
    assert.dom(GENERAL.icon('edit')).exists(); // verify password field is enabled for edit & enable button is rendered bc role hasn't been rotated
  });

  test('enterprise: it should not show edit button for password when role has been rotated', async function (assert) {
    this.version.type = 'enterprise';
    this.server.post('/sys/capabilities-self', capabilitiesStub('database/static-creds/my-role', ['update']));

    this.modelStatic.last_vault_rotation = '2025-04-21T12:51:59.063124-04:00'; // Setting a sample rotation time here to simulate what returns from BE after rotation
    await render(hbs`<DatabaseRoleEdit @model={{this.modelStatic}} @mode="edit"/>`);
    assert.dom(GENERAL.icon('edit')).doesNotExist(); // verify password field is disabled for edit & enable button isn't rendered bc role has already been rotated
  });

  test('enterprise: it should successfully create user that does rotate immediately & verify warning modal pops up', async function (assert) {
    this.version.type = 'enterprise';
    this.server.post('/sys/capabilities-self', capabilitiesStub('database/static-creds/my-role', ['create']));

    await render(hbs`<DatabaseRoleEdit @model={{this.modelStatic}} @mode="create"/>`);
    await click(GENERAL.submitButton);

    assert.dom('[data-test-issuer-warning]').exists(); // check if warning modal shows after clicking save
    await click('[data-test-issuer-save]'); // click continue button on modal

    await render(hbs`<DatabaseRoleEdit @model={{this.modelStatic}} @mode="show"/>`);
    assert.dom(GENERAL.infoRowValue('Rotate immediately')).containsText('Yes');
  });

  test('it should show Get credentials button when a user has the correct policy', async function (assert) {
    this.server.post('/sys/capabilities-self', capabilitiesStub('database/static-creds/my-role', ['read']));
    await render(hbs`<DatabaseRoleEdit @model={{this.modelStatic}} @mode="show"/>`);
    assert.dom(GENERAL.button('static')).exists('Get credentials button exists');
  });

  test('it should show Generate credentials button when a user has the correct policy', async function (assert) {
    this.server.post('/sys/capabilities-self', capabilitiesStub('database/creds/my-role', ['read']));
    await render(hbs`<DatabaseRoleEdit @model={{this.modelDynamic}} @mode="show"/>`);
    assert.dom(GENERAL.button('dynamic')).exists('Generate credentials button exists');
  });
});
