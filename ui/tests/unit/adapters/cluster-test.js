/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { resolve } from 'rsvp';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Adapter | cluster', function (hooks) {
  setupTest(hooks);

  test('cluster api urls', function (assert) {
    let url, method, options;
    const adapter = this.owner.factoryFor('adapter:cluster').create({
      ajax: (...args) => {
        [url, method, options] = args;
        return resolve();
      },
    });
    adapter.health();
    assert.strictEqual(url, '/v1/sys/health', 'health url OK');
    assert.deepEqual(
      {
        standbycode: 200,
        sealedcode: 200,
        uninitcode: 200,
        drsecondarycode: 200,
        performancestandbycode: 200,
      },
      options.data,
      'health data params OK'
    );
    assert.strictEqual(method, 'GET', 'health method OK');

    adapter.sealStatus();
    assert.strictEqual(url, '/v1/sys/seal-status', 'health url OK');
    assert.strictEqual(method, 'GET', 'seal-status method OK');

    let data = { someData: 1 };
    adapter.unseal(data);
    assert.strictEqual(url, '/v1/sys/unseal', 'unseal url OK');
    assert.strictEqual(method, 'PUT', 'unseal method OK');
    assert.deepEqual({ data, unauthenticated: true }, options, 'unseal options OK');

    adapter.initCluster(data);
    assert.strictEqual(url, '/v1/sys/init', 'init url OK');
    assert.strictEqual(method, 'PUT', 'init method OK');
    assert.deepEqual({ data, unauthenticated: true }, options, 'init options OK');

    data = { token: 'token', password: 'password', username: 'username' };

    adapter.authenticate({ backend: 'token', data });
    assert.strictEqual(url, '/v1/auth/token/lookup-self', 'auth:token url OK');
    assert.strictEqual(method, 'GET', 'auth:token method OK');
    assert.deepEqual(
      { headers: { 'X-Vault-Token': 'token' }, unauthenticated: true },
      options,
      'auth:token options OK'
    );

    adapter.authenticate({ backend: 'github', data });
    assert.strictEqual(url, '/v1/auth/github/login', 'auth:github url OK');
    assert.strictEqual(method, 'POST', 'auth:github method OK');
    assert.deepEqual(
      { data: { password: 'password', token: 'token' }, unauthenticated: true },
      options,
      'auth:github options OK'
    );

    data = { jwt: 'token', role: 'test' };
    adapter.authenticate({ backend: 'jwt', data });
    assert.strictEqual(url, '/v1/auth/jwt/login', 'auth:jwt url OK');
    assert.strictEqual(method, 'POST', 'auth:jwt method OK');
    assert.deepEqual(
      { data: { jwt: 'token', role: 'test' }, unauthenticated: true },
      options,
      'auth:jwt options OK'
    );

    data = { jwt: 'token', role: 'test', path: 'oidc' };
    adapter.authenticate({ backend: 'jwt', data });
    assert.strictEqual(url, '/v1/auth/oidc/login', 'auth:jwt custom mount path, url OK');

    data = { token: 'token', password: 'password', username: 'username', path: 'path' };

    adapter.authenticate({ backend: 'token', data });
    assert.strictEqual(url, '/v1/auth/token/lookup-self', 'auth:token url with path OK');

    adapter.authenticate({ backend: 'github', data });
    assert.strictEqual(url, '/v1/auth/path/login', 'auth:github with path url OK');

    data = { password: 'password', username: 'username' };

    adapter.authenticate({ backend: 'userpass', data });
    assert.strictEqual(url, '/v1/auth/userpass/login/username', 'auth:userpass url OK');
    assert.strictEqual(method, 'POST', 'auth:userpass method OK');
    assert.deepEqual(
      { data: { password: 'password' }, unauthenticated: true },
      options,
      'auth:userpass options OK'
    );

    adapter.authenticate({ backend: 'radius', data });
    assert.strictEqual(url, '/v1/auth/radius/login/username', 'auth:RADIUS url OK');
    assert.strictEqual(method, 'POST', 'auth:RADIUS method OK');
    assert.deepEqual(
      { data: { password: 'password' }, unauthenticated: true },
      options,
      'auth:RADIUS options OK'
    );

    adapter.authenticate({ backend: 'LDAP', data });
    assert.strictEqual(url, '/v1/auth/ldap/login/username', 'ldap:userpass url OK');
    assert.strictEqual(method, 'POST', 'ldap:userpass method OK');
    assert.deepEqual(
      { data: { password: 'password' }, unauthenticated: true },
      options,
      'ldap:userpass options OK'
    );

    data = { password: 'password', username: 'username', nonce: 'uuid' };
    adapter.authenticate({ backend: 'okta', data });
    assert.strictEqual(url, '/v1/auth/okta/login/username', 'okta:userpass url OK');
    assert.strictEqual(method, 'POST', 'ldap:userpass method OK');
    assert.deepEqual(
      { data: { password: 'password', nonce: 'uuid' }, unauthenticated: true },
      options,
      'okta:userpass options OK'
    );

    // use a custom mount path
    data = { password: 'password', username: 'username', path: 'path' };

    adapter.authenticate({ backend: 'userpass', data });
    assert.strictEqual(url, '/v1/auth/path/login/username', 'auth:userpass with path url OK');

    adapter.authenticate({ backend: 'LDAP', data });
    assert.strictEqual(url, '/v1/auth/path/login/username', 'auth:LDAP with path url OK');

    data = { password: 'password', username: 'username', path: 'path', nonce: 'uuid' };
    adapter.authenticate({ backend: 'Okta', data });
    assert.strictEqual(url, '/v1/auth/path/login/username', 'auth:Okta with path url OK');
  });

  test('cluster replication api urls', function (assert) {
    let url, method, options;
    const adapter = this.owner.factoryFor('adapter:cluster').create({
      ajax: (...args) => {
        [url, method, options] = args;
        return resolve();
      },
    });

    adapter.replicationStatus();
    assert.strictEqual(url, '/v1/sys/replication/status', 'replication:status url OK');
    assert.strictEqual(method, 'GET', 'replication:status method OK');
    assert.deepEqual({ unauthenticated: true }, options, 'replication:status options OK');

    adapter.replicationAction('recover', 'dr');
    assert.strictEqual(url, '/v1/sys/replication/recover', 'replication: recover url OK');
    assert.strictEqual(method, 'POST', 'replication:recover method OK');

    adapter.replicationAction('reindex', 'dr');
    assert.strictEqual(url, '/v1/sys/replication/reindex', 'replication: reindex url OK');
    assert.strictEqual(method, 'POST', 'replication:reindex method OK');

    adapter.replicationAction('enable', 'dr', 'primary');
    assert.strictEqual(url, '/v1/sys/replication/dr/primary/enable', 'replication:dr primary:enable url OK');
    assert.strictEqual(method, 'POST', 'replication:primary:enable method OK');
    adapter.replicationAction('enable', 'performance', 'primary');
    assert.strictEqual(
      url,
      '/v1/sys/replication/performance/primary/enable',
      'replication:performance primary:enable url OK'
    );

    adapter.replicationAction('enable', 'dr', 'secondary');
    assert.strictEqual(
      url,
      '/v1/sys/replication/dr/secondary/enable',
      'replication:dr secondary:enable url OK'
    );
    assert.strictEqual(method, 'POST', 'replication:secondary:enable method OK');
    adapter.replicationAction('enable', 'performance', 'secondary');
    assert.strictEqual(
      url,
      '/v1/sys/replication/performance/secondary/enable',
      'replication:performance secondary:enable url OK'
    );

    adapter.replicationAction('disable', 'dr', 'primary');
    assert.strictEqual(
      url,
      '/v1/sys/replication/dr/primary/disable',
      'replication:dr primary:disable url OK'
    );
    assert.strictEqual(method, 'POST', 'replication:primary:disable method OK');
    adapter.replicationAction('disable', 'performance', 'primary');
    assert.strictEqual(
      url,
      '/v1/sys/replication/performance/primary/disable',
      'replication:performance primary:disable url OK'
    );

    adapter.replicationAction('disable', 'dr', 'secondary');
    assert.strictEqual(
      url,
      '/v1/sys/replication/dr/secondary/disable',
      'replication: drsecondary:disable url OK'
    );
    assert.strictEqual(method, 'POST', 'replication:secondary:disable method OK');
    adapter.replicationAction('disable', 'performance', 'secondary');
    assert.strictEqual(
      url,
      '/v1/sys/replication/performance/secondary/disable',
      'replication: performance:disable url OK'
    );

    adapter.replicationAction('demote', 'dr', 'primary');
    assert.strictEqual(url, '/v1/sys/replication/dr/primary/demote', 'replication: dr primary:demote url OK');
    assert.strictEqual(method, 'POST', 'replication:primary:demote method OK');
    adapter.replicationAction('demote', 'performance', 'primary');
    assert.strictEqual(
      url,
      '/v1/sys/replication/performance/primary/demote',
      'replication: performance primary:demote url OK'
    );

    adapter.replicationAction('promote', 'performance', 'secondary');
    assert.strictEqual(method, 'POST', 'replication:secondary:promote method OK');
    assert.strictEqual(
      url,
      '/v1/sys/replication/performance/secondary/promote',
      'replication:performance secondary:promote url OK'
    );

    adapter.replicationDrPromote();
    assert.strictEqual(
      url,
      '/v1/sys/replication/dr/secondary/promote',
      'replication:dr secondary:promote url OK'
    );
    assert.strictEqual(method, 'PUT', 'replication:dr secondary:promote method OK');
    adapter.replicationDrPromote({}, { checkStatus: true });
    assert.strictEqual(
      url,
      '/v1/sys/replication/dr/secondary/promote',
      'replication:dr secondary:promote url OK'
    );
    assert.strictEqual(method, 'GET', 'replication:dr secondary:promote method OK');
  });
});
