/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { run } from '@ember/runloop';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { TOKEN_SEPARATOR, TOKEN_PREFIX, ROOT_PREFIX } from 'vault/services/auth';
import { TOKEN_DATA } from 'vault/tests/helpers/auth/response-stubs';

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
  ...TOKEN_DATA.token,
  policies: ['root'],
  ttl: 0, // root tokens have no expiration
};

const BATCH_TOKEN_RESPONSE = {
  ...TOKEN_DATA.token,
  renewable: false,
  type: 'batch',
};

const USERPASS_BATCH_TOKEN_RESPONSE = {
  ...TOKEN_DATA.userpass,
  renewable: false,
  tokenType: 'batch',
};

module('Integration | Service | auth', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.owner.lookup('service:flash-messages').registerTypes(['warning']);
    this.store = storage();
    this.memStore = storage();
  });

  test('token authentication: root token', function (assert) {
    assert.expect(5);
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
        .authSuccess('1', ROOT_TOKEN_RESPONSE)
        .then(() => {
          const clusterTokenName = service.get('currentTokenName');
          const clusterToken = service.get('currentToken');
          const authData = service.get('authData');

          const expectedTokenName = `${TOKEN_PREFIX}${ROOT_PREFIX}${TOKEN_SEPARATOR}1`;
          assert.strictEqual(clusterToken, 'hvs.myvaultgeneratedtoken', 'token is saved properly');
          assert.strictEqual(
            `${TOKEN_PREFIX}${ROOT_PREFIX}${TOKEN_SEPARATOR}1`,
            clusterTokenName,
            'token name is saved properly'
          );
          assert.strictEqual(authData.authMethodType, 'token', 'backend is saved properly');
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
    await service.authSuccess('1', ROOT_TOKEN_RESPONSE);
    const clusterTokenName = service.get('currentTokenName');
    const clusterToken = service.get('currentToken');
    const authData = service.get('authData');

    const expectedTokenName = `${TOKEN_PREFIX}${ROOT_PREFIX}${TOKEN_SEPARATOR}1`;
    assert.strictEqual(clusterToken, 'hvs.myvaultgeneratedtoken', 'token is saved properly');
    assert.strictEqual(
      `${TOKEN_PREFIX}${ROOT_PREFIX}${TOKEN_SEPARATOR}1`,
      clusterTokenName,
      'token name is saved properly'
    );
    assert.strictEqual(authData.authMethodType, 'token', 'backend is saved properly');
    assert.ok(this.store.keys().includes(expectedTokenName), 'root token is stored in the store');
    assert.strictEqual(this.memStore.keys().length, 0, 'mem storage is empty');
  });

  test('github authentication', function (assert) {
    assert.expect(5);
    const done = assert.async();
    const service = this.owner.factoryFor('service:auth').create({
      storage: (type) => (type === 'memory' ? this.memStore : this.store),
    });

    run(() => {
      service.authSuccess('1', TOKEN_DATA.github).then(() => {
        const clusterTokenName = service.get('currentTokenName');
        const clusterToken = service.get('currentToken');
        const authData = service.get('authData');
        const expectedTokenName = `${TOKEN_PREFIX}github${TOKEN_SEPARATOR}1`;

        assert.strictEqual(TOKEN_DATA.github.token, clusterToken, 'token is saved properly');
        assert.strictEqual(expectedTokenName, clusterTokenName, 'token name is saved properly');
        assert.strictEqual(authData.authMethodType, 'github', 'backend is saved properly');
        assert.strictEqual(this.memStore.keys().length, 0, 'mem storage is empty');
        assert.ok(this.store.keys().includes(expectedTokenName), 'normal storage contains the token');
        done();
      });
    });
  });

  test('userpass authentication', function (assert) {
    assert.expect(3);
    const done = assert.async();
    const service = this.owner.factoryFor('service:auth').create({ storage: () => this.store });
    run(() => {
      service.authSuccess('1', TOKEN_DATA.userpass).then(() => {
        const clusterTokenName = service.get('currentTokenName');
        const clusterToken = service.get('currentToken');
        const authData = service.get('authData');

        assert.strictEqual(TOKEN_DATA.userpass.token, clusterToken, 'token is saved properly');
        assert.strictEqual(
          `${TOKEN_PREFIX}userpass${TOKEN_SEPARATOR}1`,
          clusterTokenName,
          'token name is saved properly'
        );
        assert.strictEqual(authData.authMethodType, 'userpass', 'backend is saved properly');
        done();
      });
    });
  });

  test('token auth expiry with non-root token', function (assert) {
    assert.expect(4);

    const done = assert.async();
    const service = this.owner.factoryFor('service:auth').create({ storage: () => this.store });
    run(() => {
      service.authSuccess('1', TOKEN_DATA.token).then(() => {
        const clusterTokenName = service.get('currentTokenName');
        const clusterToken = service.get('currentToken');
        const authData = service.get('authData');

        assert.strictEqual(clusterToken, 'hvs.myvaultgeneratedtoken', 'token is saved properly');
        assert.strictEqual(
          `${TOKEN_PREFIX}token${TOKEN_SEPARATOR}1`,
          clusterTokenName,
          'token name is saved properly'
        );
        assert.strictEqual(authData.authMethodType, 'token', 'backend is saved properly');
        assert.false(service.get('tokenExpired'), 'token is not expired');
        done();
      });
    });
  });

  module('token types', function (hooks) {
    hooks.beforeEach(function () {
      this.service = this.owner.factoryFor('service:auth').create({ storage: () => this.store });
    });

    test('batch tokens generated by token auth method', async function (assert) {
      await this.service.authSuccess('1', BATCH_TOKEN_RESPONSE);

      // exact expiration time is calculated in unit tests
      assert.notEqual(
        this.service.tokenExpirationDate,
        undefined,
        'expiration is calculated for batch tokens'
      );
    });

    test('batch tokens generated by auth methods', async function (assert) {
      await this.service.authSuccess('1', USERPASS_BATCH_TOKEN_RESPONSE);

      // exact expiration time is calculated in unit tests
      assert.notEqual(
        this.service.tokenExpirationDate,
        undefined,
        'expiration is calculated for batch tokens'
      );
    });

    test('service token authentication', async function (assert) {
      await this.service.authSuccess('1', TOKEN_DATA.userpass);

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
