/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import sinon from 'sinon';
import { ldapRoleID } from 'vault/adapters/ldap/role';

module('Unit | Adapter | ldap/role', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.adapter = this.store.adapterFor('ldap/role');
    this.path = 'role';

    this.getModel = (type, roleName) => {
      const name = roleName || 'test-role';
      this.store.pushPayload('ldap/role', {
        modelName: 'ldap/role',
        backend: 'ldap-test',
        name,
        type,
        id: ldapRoleID(type, name),
      });
      return this.store.peekRecord('ldap/role', ldapRoleID(type, name));
    };
  });

  module('happy paths', function () {
    test('it should make request to correct endpoints when listing records', async function (assert) {
      assert.expect(6);

      const assertRequest = (schema, req) => {
        assert.ok(req.queryParams.list, 'list query param sent when listing roles');
        const name = req.url.includes('static-role') ? 'static-test' : 'dynamic-test';
        return { data: { keys: [name] } };
      };

      this.server.get('/ldap-test/static-role', assertRequest);
      this.server.get('/ldap-test/role', assertRequest);

      this.models = await this.store.query('ldap/role', { backend: 'ldap-test' });

      const model = this.models[0];
      assert.strictEqual(this.models.length, 2, 'Returns responses from both endpoints');
      assert.strictEqual(model.backend, 'ldap-test', 'Backend value is set on records returned from query');
      // sorted alphabetically by name so dynamic should be first
      assert.strictEqual(model.type, 'dynamic', 'Type value is set on records returned from query');
      assert.strictEqual(model.name, 'dynamic-test', 'Name value is set on records returned from query');
    });

    test('it should conditionally trigger info level flash message for single endpoint error from query', async function (assert) {
      const flashMessages = this.owner.lookup('service:flashMessages');
      const flashSpy = sinon.spy(flashMessages, 'info');

      this.server.get('/ldap-test/static-role', () => {
        return new Response(403, {}, { errors: ['permission denied'] });
      });
      this.server.get('/ldap-test/role', () => ({ data: { keys: ['dynamic-test'] } }));

      await this.store.query('ldap/role', { backend: 'ldap-test' });
      await this.store.query(
        'ldap/role',
        { backend: 'ldap-test' },
        { adapterOptions: { showPartialError: true } }
      );

      assert.true(
        flashSpy.calledOnceWith('Error fetching roles from /v1/ldap-test/static-role: permission denied'),
        'Partial error info only displays when adapter option is passed'
      );
    });

    test('it should throw error for query when requests to both endpoints fail', async function (assert) {
      assert.expect(2);

      this.server.get('/ldap-test/:path', (schema, req) => {
        const errors = {
          'static-role': ['permission denied'],
          role: ['server error'],
        }[req.params.path];
        return new Response(req.params.path === 'static-role' ? 403 : 500, {}, { errors });
      });

      try {
        await this.store.query('ldap/role', { backend: 'ldap-test' });
      } catch (error) {
        assert.deepEqual(
          error.errors,
          ['/v1/ldap-test/static-role: permission denied', '/v1/ldap-test/role: server error'],
          'Error messages is thrown with correct payload from query.'
        );
        assert.strictEqual(
          error.message,
          'Error fetching roles:',
          'Error message is thrown with correct payload from query.'
        );
      }
    });

    test('it should make request to correct endpoints when querying record', async function (assert) {
      assert.expect(5);

      this.server.get('/ldap-test/:path/test-role', (schema, req) => {
        assert.strictEqual(
          req.params.path,
          this.path,
          'GET request made to correct endpoint when querying record'
        );
      });

      for (const type of ['dynamic', 'static']) {
        this.model = await this.store.queryRecord('ldap/role', {
          backend: 'ldap-test',
          type,
          name: 'test-role',
        });
        this.path = 'static-role';
      }

      assert.strictEqual(
        this.model.backend,
        'ldap-test',
        'Backend value is set on records returned from query'
      );
      assert.strictEqual(this.model.type, 'static', 'Type value is set on records returned from query');
      assert.strictEqual(this.model.name, 'test-role', 'Name value is set on records returned from query');
    });

    test('it should make request to correct endpoints when creating new dynamic role record', async function (assert) {
      assert.expect(1);

      this.server.post('/ldap-test/:path/:name', (schema, req) => {
        assert.strictEqual(
          req.params.path,
          this.path,
          'POST request made to correct endpoint when creating new record for a dynamic role'
        );
      });

      const getModel = (type, name) => {
        return this.store.createRecord('ldap/role', {
          backend: 'ldap-test',
          name,
          type,
        });
      };

      const model = getModel('dynamic-role', 'dynamic-role-name');
      await model.save();
    });

    test('it should make request to correct endpoints when creating new static role record', async function (assert) {
      assert.expect(1);

      this.server.post('/ldap-test/:path/:name', (schema, req) => {
        assert.strictEqual(
          req.params.path,
          this.path,
          'POST request made to correct endpoint when creating new record for a static role'
        );
      });

      const getModel = (type, name) => {
        return this.store.createRecord('ldap/role', {
          backend: 'ldap-test',
          name,
          type,
        });
      };

      const model = getModel('static-role', 'static-role-name');
      await model.save();
    });

    test('it should make request to correct endpoints when updating record', async function (assert) {
      assert.expect(2);

      this.server.post('/ldap-test/:path/test-role', (schema, req) => {
        assert.strictEqual(
          req.params.path,
          this.path,
          'POST request made to correct endpoint when updating record'
        );
      });

      for (const type of ['dynamic', 'static']) {
        const record = this.getModel(type);
        await record.save();
        this.path = 'static-role';
      }
    });

    test('it should make request to correct endpoints when deleting record', async function (assert) {
      assert.expect(2);

      this.server.delete('/ldap-test/:path/test-role', (schema, req) => {
        assert.strictEqual(
          req.params.path,
          this.path,
          'DELETE request made to correct endpoint when deleting record'
        );
      });

      for (const type of ['dynamic', 'static']) {
        const record = this.getModel(type);
        await record.destroyRecord();
        this.path = 'static-role';
      }
    });

    test('it should make request to correct endpoints when fetching credentials', async function (assert) {
      assert.expect(2);

      this.path = 'creds';

      this.server.get('/ldap-test/:path/test-role', (schema, req) => {
        assert.strictEqual(
          req.params.path,
          this.path,
          'GET request made to correct endpoint when fetching credentials'
        );
      });

      for (const type of ['dynamic', 'static']) {
        await this.adapter.fetchCredentials('ldap-test', type, 'test-role');
        this.path = 'static-cred';
      }
    });

    test('it should make request to correct endpoint when rotating static role password', async function (assert) {
      assert.expect(1);

      this.server.post('/ldap-test/rotate-role/test-role', () => {
        assert.ok('GET request made to correct endpoint when rotating static role password');
      });

      await this.adapter.rotateStaticPassword('ldap-test', 'test-role');
    });
  });

  module('hierarchical paths', function () {
    test('it should make request to correct endpoint when listing hierarchical records', async function (assert) {
      assert.expect(2);

      const staticAncestry = { path_to_role: 'static-admin/', type: 'static' };
      const dynamicAncestry = { path_to_role: 'dynamic-admin/', type: 'dynamic' };

      this.server.get(`/ldap-test/static-role/${staticAncestry.path_to_role}`, (schema, req) => {
        assert.strictEqual(
          req.queryParams.list,
          'true',
          `query request lists roles of type: ${staticAncestry.type}`
        );
        return { data: { keys: ['my-static-role'] } };
      });
      this.server.get(`/ldap-test/role/${dynamicAncestry.path_to_role}`, (schema, req) => {
        assert.strictEqual(
          req.queryParams.list,
          'true',
          `query request lists roles of type: ${dynamicAncestry.type}`
        );
        return { data: { keys: ['my-dynamic-role'] } };
      });

      await this.store.query(
        'ldap/role',
        { backend: 'ldap-test' },
        { adapterOptions: { roleAncestry: staticAncestry } }
      );
      await this.store.query(
        'ldap/role',
        { backend: 'ldap-test' },
        { adapterOptions: { roleAncestry: dynamicAncestry } }
      );
    });

    for (const type of ['dynamic', 'static']) {
      test(`it should make request to correct endpoint when deleting a role for type: ${type}`, async function (assert) {
        assert.expect(1);

        const url =
          type === 'static'
            ? '/ldap-test/static-role/admin/my-static-role'
            : '/ldap-test/role/admin/my-dynamic-role';

        this.server.delete(url, () => {
          assert.true(true, `DELETE request made to delete hierarchical role of type: ${type}`);
        });

        const record = this.getModel(type, `admin/my-${type}-role`);
        await record.destroyRecord();
      });

      test(`it should make request to correct endpoints when fetching credentials for type: ${type}`, async function (assert) {
        assert.expect(1);

        const url =
          type === 'static'
            ? '/ldap-test/static-cred/admin/my-static-role'
            : '/ldap-test/creds/admin/my-dynamic-role';

        this.server.get(url, () => {
          assert.true(true, `request made to fetch credentials for role type: ${type}`);
        });

        await this.adapter.fetchCredentials('ldap-test', type, `admin/my-${type}-role`);
      });
    }

    test('it should make request to correct endpoint when rotating static role password', async function (assert) {
      assert.expect(1);

      this.server.post('/ldap-test/rotate-role/admin/test-role', () => {
        assert.ok('GET request made to correct endpoint when rotating static role password');
      });

      await this.adapter.rotateStaticPassword('ldap-test', 'admin/test-role');
    });
  });
});
