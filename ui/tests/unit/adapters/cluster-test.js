import { resolve } from 'rsvp';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Adapter | cluster', function (hooks) {
  setupTest(hooks);

  test('cluster api urls', function (assert) {
    let url, method, options;
    let adapter = this.owner.factoryFor('adapter:cluster').create({
      ajax: (...args) => {
        [url, method, options] = args;
        return resolve();
      },
    });
    adapter.health();
    assert.equal(url, '/v1/sys/health', 'health url OK');
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
    assert.equal(method, 'GET', 'health method OK');

    adapter.sealStatus();
    assert.equal(url, '/v1/sys/seal-status', 'health url OK');
    assert.equal(method, 'GET', 'seal-status method OK');

    let data = { someData: 1 };
    adapter.unseal(data);
    assert.equal(url, '/v1/sys/unseal', 'unseal url OK');
    assert.equal(method, 'PUT', 'unseal method OK');
    assert.deepEqual({ data, unauthenticated: true }, options, 'unseal options OK');

    adapter.initCluster(data);
    assert.equal(url, '/v1/sys/init', 'init url OK');
    assert.equal(method, 'PUT', 'init method OK');
    assert.deepEqual({ data, unauthenticated: true }, options, 'init options OK');

    data = { token: 'token', password: 'password', username: 'username' };

    adapter.authenticate({ backend: 'token', data });
    assert.equal(url, '/v1/auth/token/lookup-self', 'auth:token url OK');
    assert.equal(method, 'GET', 'auth:token method OK');
    assert.deepEqual(
      { headers: { 'X-Vault-Token': 'token' }, unauthenticated: true },
      options,
      'auth:token options OK'
    );

    adapter.authenticate({ backend: 'github', data });
    assert.equal(url, '/v1/auth/github/login', 'auth:github url OK');
    assert.equal(method, 'POST', 'auth:github method OK');
    assert.deepEqual(
      { data: { password: 'password', token: 'token' }, unauthenticated: true },
      options,
      'auth:github options OK'
    );

    data = { jwt: 'token', role: 'test' };
    adapter.authenticate({ backend: 'jwt', data });
    assert.equal(url, '/v1/auth/jwt/login', 'auth:jwt url OK');
    assert.equal(method, 'POST', 'auth:jwt method OK');
    assert.deepEqual(
      { data: { jwt: 'token', role: 'test' }, unauthenticated: true },
      options,
      'auth:jwt options OK'
    );

    data = { jwt: 'token', role: 'test', path: 'oidc' };
    adapter.authenticate({ backend: 'jwt', data });
    assert.equal(url, '/v1/auth/oidc/login', 'auth:jwt custom mount path, url OK');

    data = { token: 'token', password: 'password', username: 'username', path: 'path' };

    adapter.authenticate({ backend: 'token', data });
    assert.equal(url, '/v1/auth/token/lookup-self', 'auth:token url with path OK');

    adapter.authenticate({ backend: 'github', data });
    assert.equal(url, '/v1/auth/path/login', 'auth:github with path url OK');

    data = { password: 'password', username: 'username' };

    adapter.authenticate({ backend: 'userpass', data });
    assert.equal(url, '/v1/auth/userpass/login/username', 'auth:userpass url OK');
    assert.equal(method, 'POST', 'auth:userpass method OK');
    assert.deepEqual(
      { data: { password: 'password' }, unauthenticated: true },
      options,
      'auth:userpass options OK'
    );

    adapter.authenticate({ backend: 'radius', data });
    assert.equal(url, '/v1/auth/radius/login/username', 'auth:RADIUS url OK');
    assert.equal(method, 'POST', 'auth:RADIUS method OK');
    assert.deepEqual(
      { data: { password: 'password' }, unauthenticated: true },
      options,
      'auth:RADIUS options OK'
    );

    adapter.authenticate({ backend: 'LDAP', data });
    assert.equal(url, '/v1/auth/ldap/login/username', 'ldap:userpass url OK');
    assert.equal(method, 'POST', 'ldap:userpass method OK');
    assert.deepEqual(
      { data: { password: 'password' }, unauthenticated: true },
      options,
      'ldap:userpass options OK'
    );

    adapter.authenticate({ backend: 'okta', data });
    assert.equal(url, '/v1/auth/okta/login/username', 'okta:userpass url OK');
    assert.equal(method, 'POST', 'ldap:userpass method OK');
    assert.deepEqual(
      { data: { password: 'password' }, unauthenticated: true },
      options,
      'okta:userpass options OK'
    );

    // use a custom mount path
    data = { password: 'password', username: 'username', path: 'path' };

    adapter.authenticate({ backend: 'userpass', data });
    assert.equal(url, '/v1/auth/path/login/username', 'auth:userpass with path url OK');

    adapter.authenticate({ backend: 'LDAP', data });
    assert.equal(url, '/v1/auth/path/login/username', 'auth:LDAP with path url OK');

    adapter.authenticate({ backend: 'Okta', data });
    assert.equal(url, '/v1/auth/path/login/username', 'auth:Okta with path url OK');
  });

  test('cluster replication api urls', function (assert) {
    let url, method, options;
    let adapter = this.owner.factoryFor('adapter:cluster').create({
      ajax: (...args) => {
        [url, method, options] = args;
        return resolve();
      },
    });

    adapter.replicationStatus();
    assert.equal(url, '/v1/sys/replication/status', 'replication:status url OK');
    assert.equal(method, 'GET', 'replication:status method OK');
    assert.deepEqual({ unauthenticated: true }, options, 'replication:status options OK');

    adapter.replicationAction('recover', 'dr');
    assert.equal(url, '/v1/sys/replication/recover', 'replication: recover url OK');
    assert.equal(method, 'POST', 'replication:recover method OK');

    adapter.replicationAction('reindex', 'dr');
    assert.equal(url, '/v1/sys/replication/reindex', 'replication: reindex url OK');
    assert.equal(method, 'POST', 'replication:reindex method OK');

    adapter.replicationAction('enable', 'dr', 'primary');
    assert.equal(url, '/v1/sys/replication/dr/primary/enable', 'replication:dr primary:enable url OK');
    assert.equal(method, 'POST', 'replication:primary:enable method OK');
    adapter.replicationAction('enable', 'performance', 'primary');
    assert.equal(
      url,
      '/v1/sys/replication/performance/primary/enable',
      'replication:performance primary:enable url OK'
    );

    adapter.replicationAction('enable', 'dr', 'secondary');
    assert.equal(url, '/v1/sys/replication/dr/secondary/enable', 'replication:dr secondary:enable url OK');
    assert.equal(method, 'POST', 'replication:secondary:enable method OK');
    adapter.replicationAction('enable', 'performance', 'secondary');
    assert.equal(
      url,
      '/v1/sys/replication/performance/secondary/enable',
      'replication:performance secondary:enable url OK'
    );

    adapter.replicationAction('disable', 'dr', 'primary');
    assert.equal(url, '/v1/sys/replication/dr/primary/disable', 'replication:dr primary:disable url OK');
    assert.equal(method, 'POST', 'replication:primary:disable method OK');
    adapter.replicationAction('disable', 'performance', 'primary');
    assert.equal(
      url,
      '/v1/sys/replication/performance/primary/disable',
      'replication:performance primary:disable url OK'
    );

    adapter.replicationAction('disable', 'dr', 'secondary');
    assert.equal(url, '/v1/sys/replication/dr/secondary/disable', 'replication: drsecondary:disable url OK');
    assert.equal(method, 'POST', 'replication:secondary:disable method OK');
    adapter.replicationAction('disable', 'performance', 'secondary');
    assert.equal(
      url,
      '/v1/sys/replication/performance/secondary/disable',
      'replication: performance:disable url OK'
    );

    adapter.replicationAction('demote', 'dr', 'primary');
    assert.equal(url, '/v1/sys/replication/dr/primary/demote', 'replication: dr primary:demote url OK');
    assert.equal(method, 'POST', 'replication:primary:demote method OK');
    adapter.replicationAction('demote', 'performance', 'primary');
    assert.equal(
      url,
      '/v1/sys/replication/performance/primary/demote',
      'replication: performance primary:demote url OK'
    );

    adapter.replicationAction('promote', 'performance', 'secondary');
    assert.equal(method, 'POST', 'replication:secondary:promote method OK');
    assert.equal(
      url,
      '/v1/sys/replication/performance/secondary/promote',
      'replication:performance secondary:promote url OK'
    );

    adapter.replicationDrPromote();
    assert.equal(url, '/v1/sys/replication/dr/secondary/promote', 'replication:dr secondary:promote url OK');
    assert.equal(method, 'PUT', 'replication:dr secondary:promote method OK');
    adapter.replicationDrPromote({}, { checkStatus: true });
    assert.equal(url, '/v1/sys/replication/dr/secondary/promote', 'replication:dr secondary:promote url OK');
    assert.equal(method, 'GET', 'replication:dr secondary:promote method OK');
  });
});
