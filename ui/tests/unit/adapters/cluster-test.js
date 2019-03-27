import { resolve } from 'rsvp';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Adapter | cluster', function(hooks) {
  setupTest(hooks);

  test('cluster api urls', function(assert) {
    let url, method, options;
    let adapter = this.owner.factoryFor('adapter:cluster').create({
      ajax: (...args) => {
        [url, method, options] = args;
        return resolve();
      },
    });
    adapter.health();
    assert.equal('/v1/sys/health', url, 'health url OK');
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
    assert.equal('GET', method, 'health method OK');

    adapter.sealStatus();
    assert.equal('/v1/sys/seal-status', url, 'health url OK');
    assert.equal('GET', method, 'seal-status method OK');

    let data = { someData: 1 };
    adapter.unseal(data);
    assert.equal('/v1/sys/unseal', url, 'unseal url OK');
    assert.equal('PUT', method, 'unseal method OK');
    assert.deepEqual({ data, unauthenticated: true }, options, 'unseal options OK');

    adapter.initCluster(data);
    assert.equal('/v1/sys/init', url, 'init url OK');
    assert.equal('PUT', method, 'init method OK');
    assert.deepEqual({ data, unauthenticated: true }, options, 'init options OK');

    data = { token: 'token', password: 'password', username: 'username' };

    adapter.authenticate({ backend: 'token', data });
    assert.equal('/v1/auth/token/lookup-self', url, 'auth:token url OK');
    assert.equal('GET', method, 'auth:token method OK');
    assert.deepEqual(
      { headers: { 'X-Vault-Token': 'token' }, unauthenticated: true },
      options,
      'auth:token options OK'
    );

    adapter.authenticate({ backend: 'github', data });
    assert.equal('/v1/auth/github/login', url, 'auth:github url OK');
    assert.equal('POST', method, 'auth:github method OK');
    assert.deepEqual(
      { data: { password: 'password', token: 'token' }, unauthenticated: true },
      options,
      'auth:github options OK'
    );

    data = { jwt: 'token', role: 'test' };
    adapter.authenticate({ backend: 'jwt', data });
    assert.equal('/v1/auth/jwt/login', url, 'auth:jwt url OK');
    assert.equal('POST', method, 'auth:jwt method OK');
    assert.deepEqual(
      { data: { jwt: 'token', role: 'test' }, unauthenticated: true },
      options,
      'auth:jwt options OK'
    );

    data = { jwt: 'token', role: 'test', path: 'oidc' };
    adapter.authenticate({ backend: 'jwt', data });
    assert.equal('/v1/auth/oidc/login', url, 'auth:jwt custom mount path, url OK');

    data = { token: 'token', password: 'password', username: 'username', path: 'path' };

    adapter.authenticate({ backend: 'token', data });
    assert.equal('/v1/auth/token/lookup-self', url, 'auth:token url with path OK');

    adapter.authenticate({ backend: 'github', data });
    assert.equal('/v1/auth/path/login', url, 'auth:github with path url OK');

    data = { password: 'password', username: 'username' };

    adapter.authenticate({ backend: 'userpass', data });
    assert.equal('/v1/auth/userpass/login/username', url, 'auth:userpass url OK');
    assert.equal('POST', method, 'auth:userpass method OK');
    assert.deepEqual(
      { data: { password: 'password' }, unauthenticated: true },
      options,
      'auth:userpass options OK'
    );

    adapter.authenticate({ backend: 'radius', data });
    assert.equal('/v1/auth/radius/login/username', url, 'auth:RADIUS url OK');
    assert.equal('POST', method, 'auth:RADIUS method OK');
    assert.deepEqual(
      { data: { password: 'password' }, unauthenticated: true },
      options,
      'auth:RADIUS options OK'
    );

    adapter.authenticate({ backend: 'LDAP', data });
    assert.equal('/v1/auth/ldap/login/username', url, 'ldap:userpass url OK');
    assert.equal('POST', method, 'ldap:userpass method OK');
    assert.deepEqual(
      { data: { password: 'password' }, unauthenticated: true },
      options,
      'ldap:userpass options OK'
    );

    adapter.authenticate({ backend: 'okta', data });
    assert.equal('/v1/auth/okta/login/username', url, 'okta:userpass url OK');
    assert.equal('POST', method, 'ldap:userpass method OK');
    assert.deepEqual(
      { data: { password: 'password' }, unauthenticated: true },
      options,
      'okta:userpass options OK'
    );

    // use a custom mount path
    data = { password: 'password', username: 'username', path: 'path' };

    adapter.authenticate({ backend: 'userpass', data });
    assert.equal('/v1/auth/path/login/username', url, 'auth:userpass with path url OK');

    adapter.authenticate({ backend: 'LDAP', data });
    assert.equal('/v1/auth/path/login/username', url, 'auth:LDAP with path url OK');

    adapter.authenticate({ backend: 'Okta', data });
    assert.equal('/v1/auth/path/login/username', url, 'auth:Okta with path url OK');
  });

  test('cluster replication api urls', function(assert) {
    let url, method, options;
    let adapter = this.owner.factoryFor('adapter:cluster').create({
      ajax: (...args) => {
        [url, method, options] = args;
        return resolve();
      },
    });

    adapter.replicationStatus();
    assert.equal('/v1/sys/replication/status', url, 'replication:status url OK');
    assert.equal('GET', method, 'replication:status method OK');
    assert.deepEqual({ unauthenticated: true }, options, 'replication:status options OK');

    adapter.replicationAction('recover', 'dr');
    assert.equal('/v1/sys/replication/recover', url, 'replication: recover url OK');
    assert.equal('POST', method, 'replication:recover method OK');

    adapter.replicationAction('reindex', 'dr');
    assert.equal('/v1/sys/replication/reindex', url, 'replication: reindex url OK');
    assert.equal('POST', method, 'replication:reindex method OK');

    adapter.replicationAction('enable', 'dr', 'primary');
    assert.equal('/v1/sys/replication/dr/primary/enable', url, 'replication:dr primary:enable url OK');
    assert.equal('POST', method, 'replication:primary:enable method OK');
    adapter.replicationAction('enable', 'performance', 'primary');
    assert.equal(
      '/v1/sys/replication/performance/primary/enable',
      url,
      'replication:performance primary:enable url OK'
    );

    adapter.replicationAction('enable', 'dr', 'secondary');
    assert.equal('/v1/sys/replication/dr/secondary/enable', url, 'replication:dr secondary:enable url OK');
    assert.equal('POST', method, 'replication:secondary:enable method OK');
    adapter.replicationAction('enable', 'performance', 'secondary');
    assert.equal(
      '/v1/sys/replication/performance/secondary/enable',
      url,
      'replication:performance secondary:enable url OK'
    );

    adapter.replicationAction('disable', 'dr', 'primary');
    assert.equal('/v1/sys/replication/dr/primary/disable', url, 'replication:dr primary:disable url OK');
    assert.equal('POST', method, 'replication:primary:disable method OK');
    adapter.replicationAction('disable', 'performance', 'primary');
    assert.equal(
      '/v1/sys/replication/performance/primary/disable',
      url,
      'replication:performance primary:disable url OK'
    );

    adapter.replicationAction('disable', 'dr', 'secondary');
    assert.equal('/v1/sys/replication/dr/secondary/disable', url, 'replication: drsecondary:disable url OK');
    assert.equal('POST', method, 'replication:secondary:disable method OK');
    adapter.replicationAction('disable', 'performance', 'secondary');
    assert.equal(
      '/v1/sys/replication/performance/secondary/disable',
      url,
      'replication: performance:disable url OK'
    );

    adapter.replicationAction('demote', 'dr', 'primary');
    assert.equal('/v1/sys/replication/dr/primary/demote', url, 'replication: dr primary:demote url OK');
    assert.equal('POST', method, 'replication:primary:demote method OK');
    adapter.replicationAction('demote', 'performance', 'primary');
    assert.equal(
      '/v1/sys/replication/performance/primary/demote',
      url,
      'replication: performance primary:demote url OK'
    );

    adapter.replicationAction('promote', 'performance', 'secondary');
    assert.equal('POST', method, 'replication:secondary:promote method OK');
    assert.equal(
      '/v1/sys/replication/performance/secondary/promote',
      url,
      'replication:performance secondary:promote url OK'
    );

    adapter.replicationDrPromote();
    assert.equal('/v1/sys/replication/dr/secondary/promote', url, 'replication:dr secondary:promote url OK');
    assert.equal('PUT', method, 'replication:dr secondary:promote method OK');
    adapter.replicationDrPromote({}, { checkStatus: true });
    assert.equal('/v1/sys/replication/dr/secondary/promote', url, 'replication:dr secondary:promote url OK');
    assert.equal('GET', method, 'replication:dr secondary:promote method OK');
  });
});
