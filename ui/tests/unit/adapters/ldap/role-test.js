/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import sinon from 'sinon';

module('Unit | Adapter | ldap/role', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.adapter = this.store.adapterFor('ldap/role');
    this.path = 'role';
  });

  test('it should make request to correct endpoints when listing records', async function (assert) {
    assert.expect(6);

    const assertRequest = (schema, req) => {
      assert.ok(req.queryParams.list, 'list query param sent when listing roles');
      const name = req.params.path === 'static-role' ? 'static-test' : 'dynamic-test';
      return { data: { keys: [name] } };
    };

    this.server.get('/ldap-test/static-role', assertRequest);
    this.server.get('/ldap-test/role', assertRequest);

    this.models = await this.store.query('ldap/role', { backend: 'ldap-test' });

    const model = this.models.firstObject;
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
    assert.expect(1);

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
        error,
        {
          message: 'Error fetching roles:',
          errors: ['/v1/ldap-test/static-role: permission denied', '/v1/ldap-test/role: server error'],
        },
        'Error is thrown with correct payload from query'
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

  test('it should make request to correct endpoints when creating new record', async function (assert) {
    assert.expect(2);

    this.server.post('/ldap-test/:path/test-role', (schema, req) => {
      assert.strictEqual(
        req.params.path,
        this.path,
        'POST request made to correct endpoint when creating new record'
      );
    });

    const getModel = (type) => {
      return this.store.createRecord('ldap/role', {
        backend: 'ldap-test',
        name: 'test-role',
        type,
      });
    };

    for (const type of ['dynamic', 'static']) {
      const model = getModel(type);
      await model.save();
      this.path = 'static-role';
    }
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

    this.store.pushPayload('ldap/role', {
      modelName: 'ldap/role',
      backend: 'ldap-test',
      name: 'test-role',
    });
    const record = this.store.peekRecord('ldap/role', 'test-role');

    for (const type of ['dynamic', 'static']) {
      record.type = type;
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

    const getModel = () => {
      this.store.pushPayload('ldap/role', {
        modelName: 'ldap/role',
        backend: 'ldap-test',
        name: 'test-role',
      });
      return this.store.peekRecord('ldap/role', 'test-role');
    };

    for (const type of ['dynamic', 'static']) {
      const record = getModel();
      record.type = type;
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
