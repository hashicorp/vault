/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { settled } from '@ember/test-helpers';
import sinon from 'sinon';

const PGP_KEYS = ['pgp-key-one', 'pgp-key-two', 'pgp-key-three'];

module('Unit | Controller | vault/cluster/init', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.request = null;
    // Vault accepts PUT and POST on sys/init
    const handler = (schema, req) => {
      this.request = {
        body: JSON.parse(req.requestBody),
        headers: Object.fromEntries(
          Object.entries(req.requestHeaders).map(([key, value]) => [key.toLowerCase(), value])
        ),
      };
      return { keys: ['key'], keys_base64: ['a2V5'], root_token: 'root' };
    };
    this.server.put('/sys/init', handler);
    this.server.post('/sys/init', handler);

    this.initCluster = async (sealType, payload) => {
      const controller = this.owner.lookup('controller:vault.cluster.init');
      controller.set('model', { sealType, reload: sinon.stub(), store: this.owner.lookup('service:store') });
      controller.send('initCluster', payload);
      await settled();
    };
  });

  test('it sends a JSON request without a token', async function (assert) {
    await this.initCluster('shamir', { secret_shares: '1', secret_threshold: '1' });

    assert.true(
      this.request.headers['content-type'].includes('application/json'),
      'the request is sent as JSON'
    );
    assert.notOk(this.request.headers['x-vault-token'], 'the request is unauthenticated');
  });

  test('it sends the PGP keys for a shamir seal', async function (assert) {
    await this.initCluster('shamir', {
      secret_shares: '3',
      secret_threshold: '2',
      use_pgp: true,
      pgp_keys: PGP_KEYS,
      use_pgp_for_root: true,
      root_token_pgp_key: 'root-pgp-key',
    });

    assert.deepEqual(
      this.request.body,
      {
        secret_shares: 3,
        secret_threshold: 2,
        pgp_keys: PGP_KEYS,
        root_token_pgp_key: 'root-pgp-key',
      },
      'the PGP keys are sent'
    );
  });

  test('it sends the recovery PGP keys for a cloud seal', async function (assert) {
    await this.initCluster('awskms', {
      secret_shares: '3',
      secret_threshold: '2',
      use_pgp: true,
      pgp_keys: PGP_KEYS,
      use_pgp_for_root: true,
      root_token_pgp_key: 'root-pgp-key',
    });

    assert.deepEqual(
      this.request.body,
      {
        stored_shares: 1,
        recovery_shares: 3,
        recovery_threshold: 2,
        pgp_keys: PGP_KEYS,
        recovery_pgp_keys: PGP_KEYS,
        root_token_pgp_key: 'root-pgp-key',
      },
      'the recovery PGP keys are sent'
    );
  });

  test('it omits the PGP keys when they are not selected', async function (assert) {
    await this.initCluster('shamir', {
      secret_shares: '1',
      secret_threshold: '1',
      use_pgp: false,
      pgp_keys: PGP_KEYS,
      use_pgp_for_root: false,
      root_token_pgp_key: 'root-pgp-key',
    });

    assert.deepEqual(this.request.body, { secret_shares: 1, secret_threshold: 1 }, 'no PGP keys are sent');
  });
});
