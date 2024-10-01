/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { run } from '@ember/runloop';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { TOKEN_SEPARATOR, TOKEN_PREFIX, ROOT_PREFIX } from 'vault/services/auth';
import { setupMirage } from 'ember-cli-mirage/test-support';

function storage() {
  return {
    items: {},
    getItem(key) {
      var item = this.items[key];
      return item && JSON.parse(item);
    },

    setItem(key, val) {
      return (this.items[key] = JSON.stringify(val));
    },

    removeItem(key) {
      delete this.items[key];
    },

    keys() {
      return Object.keys(this.items);
    },
  };
}

const ROOT_TOKEN_RESPONSE = {
  request_id: 'e6674d7f-c96f-d51f-4463-cc95f0ad307e',
  lease_id: '',
  renewable: false,
  lease_duration: 0,
  data: {
    accessor: '1dd25306-fdb9-0f43-8169-48ad702041b0',
    creation_time: 1477671134,
    creation_ttl: 0,
    display_name: 'root',
    explicit_max_ttl: 0,
    id: '',
    meta: null,
    num_uses: 0,
    orphan: true,
    path: 'auth/token/root',
    policies: ['root'],
    ttl: 0,
  },
  wrap_info: null,
  warnings: null,
  auth: null,
};

const TOKEN_NON_ROOT_RESPONSE = function () {
  return {
    request_id: '3ca32cd9-fd40-891d-02d5-ea23138e8642',
    lease_id: '',
    renewable: false,
    lease_duration: 0,
    data: {
      accessor: '4ef32471-a94c-79ee-c290-aeba4d63bdc9',
      creation_time: Math.floor(Date.now() / 1000),
      creation_ttl: 2764800,
      display_name: 'token',
      explicit_max_ttl: 0,
      id: '6d83e912-1b21-9df9-b51a-d201b709f3d5',
      meta: null,
      num_uses: 0,
      orphan: false,
      path: 'auth/token/create',
      policies: ['default', 'userpass'],
      renewable: true,
      ttl: 2763327,
    },
    wrap_info: null,
    warnings: null,
    auth: null,
  };
};

const USERPASS_RESPONSE = {
  request_id: '7e5e8d3d-599e-6ef7-7570-f7057fc7c53d',
  lease_id: '',
  renewable: false,
  lease_duration: 0,
  data: null,
  wrap_info: null,
  warnings: null,
  auth: {
    client_token: '5313ff81-05cb-699f-29d1-b82b4e2906dc',
    accessor: '5c5303e7-56d6-ea13-72df-d85411bd9a7d',
    policies: ['default'],
    metadata: {
      username: 'matthew',
    },
    lease_duration: 2764800,
    renewable: true,
  },
};

const GITHUB_RESPONSE = {
  request_id: '4913f9cd-a95f-d1f9-5746-4c3af4e15660',
  lease_id: '',
  renewable: false,
  lease_duration: 0,
  data: null,
  wrap_info: null,
  warnings: null,
  auth: {
    client_token: '0d39b535-598e-54d9-96e3-97493492a5f7',
    accessor: 'd8cd894f-bedf-5ce3-f1b5-98f7c6cf8ab4',
    policies: ['default'],
    metadata: {
      org: 'hashicorp',
      username: 'meirish',
    },
    lease_duration: 2764800,
    renewable: true,
  },
};

const BATCH_TOKEN_RESPONSE = {
  request_id: '60bcef62-cc20-facf-8c0d-1418d05e9a42',
  lease_id: '',
  renewable: false,
  lease_duration: 0,
  data: {
    accessor: '',
    creation_time: 1718672331,
    creation_ttl: 60,
    display_name: 'token',
    entity_id: '',
    expire_time: '2024-06-17T17:59:51-07:00',
    explicit_max_ttl: 0,
    id: 'hvb.AAAAAQIUMVkhx9rnA',
    issue_time: '2024-06-17T17:58:51-07:00',
    meta: null,
    num_uses: 0,
    orphan: false,
    path: 'auth/token/create',
    policies: ['default'],
    renewable: false,
    ttl: 45,
    type: 'batch',
  },
  wrap_info: null,
  warnings: null,
  auth: null,
  mount_type: 'token',
};

const USERPASS_BATCH_TOKEN_RESPONSE = {
  request_id: 'eb4c31a0-1745-5701-cce7-1668f5839dbf',
  lease_id: '',
  renewable: false,
  lease_duration: 0,
  data: null,
  wrap_info: null,
  warnings: null,
  auth: {
    client_token: 'hvb.AAAAAQJ0eGwP5e48S61kBRYmR',
    accessor: '',
    policies: ['default'],
    token_policies: ['default'],
    metadata: {
      username: 'bob',
    },
    lease_duration: 360,
    renewable: false,
    entity_id: 'b52f8591-02b6-828b-7f36-620afa539126',
    token_type: 'batch',
    orphan: true,
    mfa_requirement: null,
    num_uses: 0,
  },
  mount_type: '',
};

const USERPASS_SERVICE_TOKEN_RESPONSE = {
  request_id: 'e735ffad-f2fe-5d1b-14b8-90aeb9d05976',
  lease_id: '',
  renewable: false,
  lease_duration: 0,
  data: null,
  wrap_info: null,
  warnings: null,
  auth: {
    client_token: 'hvs.CAESINY6Qbs8rm',
    accessor: '9bDizzlcIHiXwEOK5mZ6gjHI',
    policies: ['default'],
    token_policies: ['default'],
    metadata: {
      username: 'bob',
    },
    lease_duration: 360,
    renewable: true,
    entity_id: 'd9a0cac8-779c-e766-716a-6f80552f0e81',
    token_type: 'service',
    orphan: true,
    mfa_requirement: null,
    num_uses: 0,
  },
  mount_type: '',
};

module('Integration | Service | auth', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.owner.lookup('service:flash-messages').registerTypes(['warning']);
    this.store = storage();
    this.memStore = storage();
    this.server.get('/auth/token/lookup-self', function (_, request) {
      const resp = { ...ROOT_TOKEN_RESPONSE };
      resp.id = request.requestHeaders['X-Vault-Token'];
      resp.data.id = request.requestHeaders['X-Vault-Token'];
      return resp;
    });
    this.server.post('/auth/userpass/login/:username', function (_, request) {
      const { username } = request.params;
      const resp = { ...USERPASS_RESPONSE };
      resp.auth.metadata.username = username;
      return resp;
    });

    this.server.post('/auth/github/login', function () {
      return { ...GITHUB_RESPONSE };
    });
  });

  test('token authentication: root token', function (assert) {
    assert.expect(6);
    const done = assert.async();
    const self = this;
    const service = this.owner.factoryFor('service:auth').create({
      storage(tokenName) {
        if (
          tokenName &&
          tokenName.indexOf(`${TOKEN_PREFIX}${ROOT_PREFIX}`) === 0 &&
          this.environment() !== 'development'
        ) {
          return self.memStore;
        } else {
          return self.store;
        }
      },
    });
    run(() => {
      service
        .authenticate({ clusterId: '1', backend: 'token', data: { token: 'test' } })
        .then(() => {
          const clusterTokenName = service.get('currentTokenName');
          const clusterToken = service.get('currentToken');
          const authData = service.get('authData');

          const expectedTokenName = `${TOKEN_PREFIX}${ROOT_PREFIX}${TOKEN_SEPARATOR}1`;
          assert.strictEqual(clusterToken, 'test', 'token is saved properly');
          assert.strictEqual(
            `${TOKEN_PREFIX}${ROOT_PREFIX}${TOKEN_SEPARATOR}1`,
            clusterTokenName,
            'token name is saved properly'
          );
          assert.strictEqual(authData.backend.type, 'token', 'backend is saved properly');
          assert.strictEqual(
            ROOT_TOKEN_RESPONSE.data.display_name,
            authData.displayName,
            'displayName is saved properly'
          );
          assert.ok(
            this.memStore.keys().includes(expectedTokenName),
            'root token is stored in the memory store'
          );
          assert.strictEqual(this.store.keys().length, 0, 'normal storage is empty');
        })
        .finally(() => {
          done();
        });
    });
  });

  test('token authentication: root token in ember development environment', async function (assert) {
    const self = this;
    const service = this.owner.factoryFor('service:auth').create({
      storage(tokenName) {
        if (
          tokenName &&
          tokenName.indexOf(`${TOKEN_PREFIX}${ROOT_PREFIX}`) === 0 &&
          this.environment() !== 'development'
        ) {
          return self.memStore;
        } else {
          return self.store;
        }
      },
      environment: () => 'development',
    });
    await service.authenticate({ clusterId: '1', backend: 'token', data: { token: 'test' } });
    const clusterTokenName = service.get('currentTokenName');
    const clusterToken = service.get('currentToken');
    const authData = service.get('authData');

    const expectedTokenName = `${TOKEN_PREFIX}${ROOT_PREFIX}${TOKEN_SEPARATOR}1`;
    assert.strictEqual(clusterToken, 'test', 'token is saved properly');
    assert.strictEqual(
      `${TOKEN_PREFIX}${ROOT_PREFIX}${TOKEN_SEPARATOR}1`,
      clusterTokenName,
      'token name is saved properly'
    );
    assert.strictEqual(authData.backend.type, 'token', 'backend is saved properly');
    assert.strictEqual(
      ROOT_TOKEN_RESPONSE.data.display_name,
      authData.displayName,
      'displayName is saved properly'
    );
    assert.ok(this.store.keys().includes(expectedTokenName), 'root token is stored in the store');
    assert.strictEqual(this.memStore.keys().length, 0, 'mem storage is empty');
  });

  test('github authentication', function (assert) {
    assert.expect(6);
    const done = assert.async();
    const service = this.owner.factoryFor('service:auth').create({
      storage: (type) => (type === 'memory' ? this.memStore : this.store),
    });

    run(() => {
      service.authenticate({ clusterId: '1', backend: 'github', data: { token: 'test' } }).then(() => {
        const clusterTokenName = service.get('currentTokenName');
        const clusterToken = service.get('currentToken');
        const authData = service.get('authData');
        const expectedTokenName = `${TOKEN_PREFIX}github${TOKEN_SEPARATOR}1`;

        assert.strictEqual(GITHUB_RESPONSE.auth.client_token, clusterToken, 'token is saved properly');
        assert.strictEqual(expectedTokenName, clusterTokenName, 'token name is saved properly');
        assert.strictEqual(authData.backend.type, 'github', 'backend is saved properly');
        assert.strictEqual(
          GITHUB_RESPONSE.auth.metadata.org + '/' + GITHUB_RESPONSE.auth.metadata.username,
          authData.displayName,
          'displayName is saved properly'
        );
        assert.strictEqual(this.memStore.keys().length, 0, 'mem storage is empty');
        assert.ok(this.store.keys().includes(expectedTokenName), 'normal storage contains the token');
        done();
      });
    });
  });

  test('userpass authentication', function (assert) {
    assert.expect(4);
    const done = assert.async();
    const service = this.owner.factoryFor('service:auth').create({ storage: () => this.store });
    run(() => {
      service
        .authenticate({
          clusterId: '1',
          backend: 'userpass',
          data: { username: USERPASS_RESPONSE.auth.metadata.username, password: 'passoword' },
        })
        .then(() => {
          const clusterTokenName = service.get('currentTokenName');
          const clusterToken = service.get('currentToken');
          const authData = service.get('authData');

          assert.strictEqual(USERPASS_RESPONSE.auth.client_token, clusterToken, 'token is saved properly');
          assert.strictEqual(
            `${TOKEN_PREFIX}userpass${TOKEN_SEPARATOR}1`,
            clusterTokenName,
            'token name is saved properly'
          );
          assert.strictEqual(authData.backend.type, 'userpass', 'backend is saved properly');
          assert.strictEqual(
            USERPASS_RESPONSE.auth.metadata.username,
            authData.displayName,
            'displayName is saved properly'
          );
          done();
        });
    });
  });

  test('token auth expiry with non-root token', function (assert) {
    assert.expect(5);
    const tokenResp = TOKEN_NON_ROOT_RESPONSE();
    this.server.get('/auth/token/lookup-self', function (_, request) {
      const resp = { ...tokenResp };
      resp.id = request.requestHeaders['X-Vault-Token'];
      resp.data.id = request.requestHeaders['X-Vault-Token'];
      return resp;
    });

    const done = assert.async();
    const service = this.owner.factoryFor('service:auth').create({ storage: () => this.store });
    run(() => {
      service.authenticate({ clusterId: '1', backend: 'token', data: { token: 'test' } }).then(() => {
        const clusterTokenName = service.get('currentTokenName');
        const clusterToken = service.get('currentToken');
        const authData = service.get('authData');

        assert.strictEqual(clusterToken, 'test', 'token is saved properly');
        assert.strictEqual(
          `${TOKEN_PREFIX}token${TOKEN_SEPARATOR}1`,
          clusterTokenName,
          'token name is saved properly'
        );
        assert.strictEqual(authData.backend.type, 'token', 'backend is saved properly');
        assert.strictEqual(
          authData.displayName,
          tokenResp.data.display_name,
          'displayName is saved properly'
        );
        assert.false(service.get('tokenExpired'), 'token is not expired');
        done();
      });
    });
  });

  module('token types', function (hooks) {
    hooks.beforeEach(function () {
      this.server.post('/auth/userpass/login/:username', (_, request) => {
        const { username } = request.params;
        const resp =
          username === 'batch'
            ? { ...USERPASS_BATCH_TOKEN_RESPONSE }
            : { ...USERPASS_SERVICE_TOKEN_RESPONSE };
        resp.auth.metadata.username = username;
        return resp;
      });

      this.service = this.owner.factoryFor('service:auth').create({ storage: () => this.store });
    });

    module('batch tokens', function () {
      test('batch tokens generated by token auth method', async function (assert) {
        this.server.get('/auth/token/lookup-self', () => {
          return { ...BATCH_TOKEN_RESPONSE };
        });

        await this.service.authenticate({
          clusterId: '1',
          backend: 'token',
          data: { token: 'test' },
        });

        // exact expiration time is calculated in unit tests
        assert.notEqual(
          this.service.tokenExpirationDate,
          undefined,
          'expiration is calculated for batch tokens'
        );
      });

      test('batch tokens generated by auth methods', async function (assert) {
        await this.service.authenticate({
          clusterId: '1',
          backend: 'userpass',
          data: { username: 'batch', password: 'password' },
        });

        // exact expiration time is calculated in unit tests
        assert.notEqual(
          this.service.tokenExpirationDate,
          undefined,
          'expiration is calculated for batch tokens'
        );
      });
    });

    test('service token authentication', async function (assert) {
      await this.service.authenticate({
        clusterId: '1',
        backend: 'userpass',
        data: { username: 'service', password: 'password' },
      });

      // exact expiration time is calculated in unit tests
      assert.notEqual(
        this.service.tokenExpirationDate,
        undefined,
        'expiration is calculated for service tokens'
      );
      assert.false(this.service.allowExpiration, 'allowExpiration is false for service tokens');
    });
  });
});
