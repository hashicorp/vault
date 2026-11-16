/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import DatabaseRoleForm, {
  getRoleSettingAttrs,
  getRoleShowFields,
  normalizeDatabaseRole,
  ROLE_DEFAULTS,
} from 'vault/forms/database/role';

module('Unit | Form | database/role', function () {
  module('toJSON', function () {
    /**
     * A new dynamic role sends its ttls, the connection as db_name and its statements,
     * and never sends static-only attributes.
     */
    test('it serializes a new dynamic role', function (assert) {
      const form = new DatabaseRoleForm(
        {
          ...ROLE_DEFAULTS,
          name: 'dyn-role',
          backend: 'database',
          database: ['my-db'],
          type: 'dynamic',
          creation_statements: ['CREATE ROLE "{{name}}"'],
          username: 'ignored',
        },
        { isNew: true }
      );
      const { isValid, data } = form.toJSON();
      assert.true(isValid, 'form is valid');
      assert.propEqual(
        data,
        {
          default_ttl: '1h',
          max_ttl: '24h',
          db_name: 'my-db',
          creation_statements: ['CREATE ROLE "{{name}}"'],
        },
        'payload only contains dynamic attributes'
      );
    });

    /**
     * A new static role sends the username, rotation period and the rotate immediately flag
     * but not dynamic attributes.
     */
    test('it serializes a new static role', function (assert) {
      const form = new DatabaseRoleForm(
        {
          ...ROLE_DEFAULTS,
          name: 'static-role',
          backend: 'database',
          database: ['my-db'],
          type: 'static',
          username: 'vault-user',
          skip_import_rotation: true,
        },
        { isNew: true }
      );
      const { isValid, data } = form.toJSON();
      assert.true(isValid, 'form is valid');
      assert.propEqual(
        data,
        {
          username: 'vault-user',
          rotation_period: '24h',
          skip_import_rotation: true,
          db_name: 'my-db',
        },
        'payload only contains static attributes'
      );
    });

    /**
     * MongoDB roles use a json string for the singular statement, which the API expects as a list.
     */
    test('it transposes a singular statement into a list', function (assert) {
      const form = new DatabaseRoleForm(
        {
          name: 'mongo-role',
          backend: 'database',
          database: ['mongo'],
          type: 'dynamic',
          creation_statement: '{"db":"admin"}',
          revocation_statement: '{"db":"admin"}',
        },
        { isNew: true }
      );
      const { data } = form.toJSON();
      assert.propEqual(data.creation_statements, ['{"db":"admin"}'], 'creation statement is a list');
      assert.propEqual(data.revocation_statements, ['{"db":"admin"}'], 'revocation statement is a list');
    });

    /**
     * On update only changed attributes are sent. The static username is always sent
     * because the API requires it, and db_name is never sent since it cannot be changed.
     */
    test('it only sends changed attributes on update of a static role', function (assert) {
      const form = new DatabaseRoleForm(
        {
          name: 'static-role',
          backend: 'database',
          database: ['my-db'],
          type: 'static',
          username: 'vault-user',
          rotation_period: '24h',
          skip_import_rotation: true,
        },
        { isNew: false }
      );
      form.rotation_period = '48h';
      const { data } = form.toJSON();
      assert.propEqual(data, { username: 'vault-user', rotation_period: '48h' }, 'only changes are sent');
    });

    test('it only sends changed attributes on update of a dynamic role', function (assert) {
      const form = new DatabaseRoleForm(
        {
          name: 'dyn-role',
          backend: 'database',
          database: ['my-db'],
          type: 'dynamic',
          default_ttl: '1h',
          max_ttl: '24h',
          creation_statements: ['CREATE'],
        },
        { isNew: false }
      );
      form.max_ttl = '48h';
      const { data } = form.toJSON();
      assert.propEqual(data, { max_ttl: '48h' }, 'only the changed ttl is sent');
    });

    /**
     * Vault keeps statements that are omitted from an update, so clearing one must send an empty list.
     */
    test('it sends an empty list when an existing statement is cleared on update', function (assert) {
      const form = new DatabaseRoleForm(
        normalizeDatabaseRole(
          {
            db_name: 'my-db',
            creation_statements: ['CREATE'],
            revocation_statements: ['REVOKE'],
            renew_statements: ['RENEW'],
          },
          'dyn-role',
          'database',
          'dynamic'
        ),
        { isNew: false }
      );
      form.revocation_statements = [];
      form.creation_statement = '';
      const { data } = form.toJSON();
      assert.propEqual(
        data,
        { creation_statements: [], revocation_statements: [] },
        'cleared statements are sent as empty lists and unchanged ones are omitted'
      );
    });

    /**
     * An empty password would overwrite the stored one, so it is omitted unless one was entered.
     */
    test('it omits an empty password and sends an entered one', function (assert) {
      const form = new DatabaseRoleForm(
        {
          name: 'static-role',
          backend: 'database',
          database: ['my-db'],
          type: 'static',
          username: 'vault-user',
          password: '',
        },
        { isNew: false }
      );
      assert.false('password' in form.toJSON().data, 'empty password is not sent');
      form.password = 'secret';
      assert.strictEqual(form.toJSON().data.password, 'secret', 'entered password is sent');
    });
  });

  module('validations', function () {
    test('it requires name, database and type', function (assert) {
      const { isValid, state } = new DatabaseRoleForm({ backend: 'database', database: [] }).toJSON();
      assert.false(isValid, 'form is invalid');
      assert.false(state.name.isValid, 'name is required');
      assert.false(state.database.isValid, 'database is required');
      assert.false(state.type.isValid, 'type is required');
    });

    test('it requires a username for static roles only', function (assert) {
      const base = { name: 'r', backend: 'database', database: ['db'] };
      const staticResult = new DatabaseRoleForm({ ...base, type: 'static' }).toJSON();
      assert.false(staticResult.isValid, 'static role without username is invalid');
      assert.strictEqual(staticResult.state.username.errors[0], 'Username is required.', 'error message');
      const dynamicResult = new DatabaseRoleForm({ ...base, type: 'dynamic' }).toJSON();
      assert.true(dynamicResult.isValid, 'dynamic role does not need a username');
    });
  });

  module('helpers', function () {
    test('normalizeDatabaseRole flattens the read response', function (assert) {
      const role = normalizeDatabaseRole(
        { db_name: 'my-db', creation_statements: ['CREATE'], default_ttl: 3600 },
        'my-role',
        'database',
        'dynamic'
      );
      assert.propEqual(
        role,
        {
          creation_statements: ['CREATE'],
          default_ttl: 3600,
          name: 'my-role',
          backend: 'database',
          type: 'dynamic',
          database: ['my-db'],
          creation_statement: 'CREATE',
          revocation_statement: '',
          idPrefix: 'role/',
        },
        'role is flattened'
      );
    });

    test('normalizeDatabaseRole handles an empty response', function (assert) {
      const role = normalizeDatabaseRole(undefined, 'r', 'database', 'static');
      assert.propEqual(role.database, [], 'database is empty');
      assert.strictEqual(role.type, 'static', 'type is set');
    });

    test('getRoleSettingAttrs removes enterprise fields on community edition', function (assert) {
      const ce = getRoleSettingAttrs(false).map((f) => f.name);
      const ent = getRoleSettingAttrs(true).map((f) => f.name);
      assert.false(ce.includes('skip_import_rotation'), 'ce has no skip_import_rotation');
      assert.false(ce.includes('password'), 'ce has no password');
      assert.true(ent.includes('skip_import_rotation'), 'ent has skip_import_rotation');
      assert.true(ent.includes('password'), 'ent has password');
    });

    test('getRoleShowFields omits revocation statements for elasticsearch', function (assert) {
      const names = (db) => getRoleShowFields('dynamic', [db]).map((f) => f.name);
      assert.true(names('postgres').includes('revocation_statements'), 'postgres shows revocation');
      assert.false(names('elasticsearch').includes('revocation_statements'), 'elasticsearch does not');
    });
  });
});
