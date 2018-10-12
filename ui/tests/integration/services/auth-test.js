import { run } from '@ember/runloop';
import { copy } from '@ember/object/internals';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { TOKEN_SEPARATOR, TOKEN_PREFIX, ROOT_PREFIX } from 'vault/services/auth';
import Pretender from 'pretender';

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

let ROOT_TOKEN_RESPONSE = {
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

let TOKEN_NON_ROOT_RESPONSE = function() {
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

let USERPASS_RESPONSE = {
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

let GITHUB_RESPONSE = {
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

module('Integration | Service | auth', function(hooks) {
  setupTest(hooks);

  hooks.beforeEach(function() {
    this.owner.lookup('service:flash-messages').registerTypes(['warning']);
    this.store = storage();
    this.memStore = storage();
    this.server = new Pretender(function() {
      this.get('/v1/auth/token/lookup-self', function(request) {
        let resp = copy(ROOT_TOKEN_RESPONSE, true);
        resp.id = request.requestHeaders['X-Vault-Token'];
        resp.data.id = request.requestHeaders['X-Vault-Token'];
        return [200, {}, resp];
      });
      this.post('/v1/auth/userpass/login/:username', function(request) {
        const { username } = request.params;
        let resp = copy(USERPASS_RESPONSE, true);
        resp.auth.metadata.username = username;
        return [200, {}, resp];
      });

      this.post('/v1/auth/github/login', function() {
        let resp = copy(GITHUB_RESPONSE, true);
        return [200, {}, resp];
      });
    });

    this.server.prepareBody = function(body) {
      return body ? JSON.stringify(body) : '{"error": "not found"}';
    };

    this.server.prepareHeaders = function(headers) {
      headers['content-type'] = 'application/javascript';
      return headers;
    };
  });

  hooks.afterEach(function() {
    this.server.shutdown();
  });

  test('token authentication: root token', function(assert) {
    let done = assert.async();
    let self = this;
    let service = this.owner.factoryFor('service:auth').create({
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
      service.authenticate({ clusterId: '1', backend: 'token', data: { token: 'test' } }).then(() => {
        const clusterTokenName = service.get('currentTokenName');
        const clusterToken = service.get('currentToken');
        const authData = service.get('authData');

        const expectedTokenName = `${TOKEN_PREFIX}${ROOT_PREFIX}${TOKEN_SEPARATOR}1`;
        assert.equal('test', clusterToken, 'token is saved properly');
        assert.equal(
          `${TOKEN_PREFIX}${ROOT_PREFIX}${TOKEN_SEPARATOR}1`,
          clusterTokenName,
          'token name is saved properly'
        );
        assert.equal('token', authData.backend.type, 'backend is saved properly');
        assert.equal(
          ROOT_TOKEN_RESPONSE.data.display_name,
          authData.displayName,
          'displayName is saved properly'
        );
        assert.ok(
          this.memStore.keys().includes(expectedTokenName),
          'root token is stored in the memory store'
        );
        assert.equal(this.store.keys().length, 0, 'normal storage is empty');
        done();
      });
    });
  });

  test('token authentication: root token in ember development environment', function(assert) {
    let done = assert.async();
    let self = this;
    let service = this.owner.factoryFor('service:auth').create({
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
    run(() => {
      service.authenticate({ clusterId: '1', backend: 'token', data: { token: 'test' } }).then(() => {
        const clusterTokenName = service.get('currentTokenName');
        const clusterToken = service.get('currentToken');
        const authData = service.get('authData');

        const expectedTokenName = `${TOKEN_PREFIX}${ROOT_PREFIX}${TOKEN_SEPARATOR}1`;
        assert.equal('test', clusterToken, 'token is saved properly');
        assert.equal(
          `${TOKEN_PREFIX}${ROOT_PREFIX}${TOKEN_SEPARATOR}1`,
          clusterTokenName,
          'token name is saved properly'
        );
        assert.equal('token', authData.backend.type, 'backend is saved properly');
        assert.equal(
          ROOT_TOKEN_RESPONSE.data.display_name,
          authData.displayName,
          'displayName is saved properly'
        );
        assert.ok(this.store.keys().includes(expectedTokenName), 'root token is stored in the store');
        assert.equal(this.memStore.keys().length, 0, 'mem storage is empty');
        done();
      });
    });
  });

  test('github authentication', function(assert) {
    let done = assert.async();
    let service = this.owner.factoryFor('service:auth').create({
      storage: type => (type === 'memory' ? this.memStore : this.store),
    });

    run(() => {
      service.authenticate({ clusterId: '1', backend: 'github', data: { token: 'test' } }).then(() => {
        const clusterTokenName = service.get('currentTokenName');
        const clusterToken = service.get('currentToken');
        const authData = service.get('authData');
        const expectedTokenName = `${TOKEN_PREFIX}github${TOKEN_SEPARATOR}1`;

        assert.equal(GITHUB_RESPONSE.auth.client_token, clusterToken, 'token is saved properly');
        assert.equal(expectedTokenName, clusterTokenName, 'token name is saved properly');
        assert.equal('github', authData.backend.type, 'backend is saved properly');
        assert.equal(
          GITHUB_RESPONSE.auth.metadata.org + '/' + GITHUB_RESPONSE.auth.metadata.username,
          authData.displayName,
          'displayName is saved properly'
        );
        assert.equal(this.memStore.keys().length, 0, 'mem storage is empty');
        assert.ok(this.store.keys().includes(expectedTokenName), 'normal storage contains the token');
        done();
      });
    });
  });

  test('userpass authentication', function(assert) {
    let done = assert.async();
    let service = this.owner.factoryFor('service:auth').create({ storage: () => this.store });
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

          assert.equal(USERPASS_RESPONSE.auth.client_token, clusterToken, 'token is saved properly');
          assert.equal(
            `${TOKEN_PREFIX}userpass${TOKEN_SEPARATOR}1`,
            clusterTokenName,
            'token name is saved properly'
          );
          assert.equal('userpass', authData.backend.type, 'backend is saved properly');
          assert.equal(
            USERPASS_RESPONSE.auth.metadata.username,
            authData.displayName,
            'displayName is saved properly'
          );
          done();
        });
    });
  });

  test('token auth expiry with non-root token', function(assert) {
    const tokenResp = TOKEN_NON_ROOT_RESPONSE();
    this.server.map(function() {
      this.get('/v1/auth/token/lookup-self', function(request) {
        let resp = copy(tokenResp, true);
        resp.id = request.requestHeaders['X-Vault-Token'];
        resp.data.id = request.requestHeaders['X-Vault-Token'];
        return [200, {}, resp];
      });
    });

    let done = assert.async();
    let service = this.owner.factoryFor('service:auth').create({ storage: () => this.store });
    run(() => {
      service.authenticate({ clusterId: '1', backend: 'token', data: { token: 'test' } }).then(() => {
        const clusterTokenName = service.get('currentTokenName');
        const clusterToken = service.get('currentToken');
        const authData = service.get('authData');

        assert.equal('test', clusterToken, 'token is saved properly');
        assert.equal(
          `${TOKEN_PREFIX}token${TOKEN_SEPARATOR}1`,
          clusterTokenName,
          'token name is saved properly'
        );
        assert.equal(authData.backend.type, 'token', 'backend is saved properly');
        assert.equal(authData.displayName, tokenResp.data.display_name, 'displayName is saved properly');
        assert.equal(service.get('tokenExpired'), false, 'token is not expired');
        done();
      });
    });
  });
});
