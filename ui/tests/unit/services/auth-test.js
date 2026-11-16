/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';

module('Unit | Service | auth', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.service = this.owner.lookup('service:auth');
  });

  module('#calculateExpiration', function () {
    test('with a non-zero ttl value', function (assert) {
      const now = Date.now();
      const ttl = 30;
      const expireTime = null;
      const calculatedExpiry = this.service.calculateExpiration({ now, ttl, expireTime });

      assert.strictEqual(calculatedExpiry.ttl, 30, 'returns the ttl');
      assert.strictEqual(
        calculatedExpiry.tokenExpirationEpoch,
        now + ttl * 1e3,
        'calculates expiration from ttl as epoch timestamp'
      );
    });

    test('with a zero ttl value', function (assert) {
      const now = Date.now();
      const ttl = 0;
      const expireTime = null;
      const calculatedExpiry = this.service.calculateExpiration({ now, ttl, expireTime });

      assert.strictEqual(calculatedExpiry.ttl, null, 'returns `null` for the ttl');
      assert.strictEqual(calculatedExpiry.tokenExpirationEpoch, null, 'tokenExpirationEpoch is null');
    });

    test('#calculateExpiration w/ expireTime', function (assert) {
      const now = Date.now();
      const expirationString = '2024-06-13T09:10:21-07:00';
      const expectedExpirationEpoch = new Date(expirationString).getTime();

      const calculatedExpiry = this.service.calculateExpiration({
        now,
        ttl: 30,
        expireTime: '2024-06-13T09:10:21-07:00',
      });

      assert.strictEqual(calculatedExpiry.ttl, 30, 'returns ttl');
      assert.strictEqual(
        calculatedExpiry.tokenExpirationEpoch,
        expectedExpirationEpoch,
        'calculates expiration from expireTime'
      );
    });
  });

  module('#setExpirationSettings', function () {
    test('#setExpirationSettings for a renewable token', function (assert) {
      const now = Date.now();
      const renewable = true;

      this.service.setExpirationSettings(renewable, now);

      assert.false(this.service.allowExpiration, 'sets allowExpiration to false');
      assert.strictEqual(this.service.expirationCalcTS, now, 'sets expirationCalcTS to now');
    });

    test('#setExpirationSettings for a non-renewable token', function (assert) {
      const now = Date.now();
      const renewable = false;

      this.service.setExpirationSettings(renewable, now);

      assert.true(this.service.allowExpiration, 'sets allowExpiration to true');
      assert.strictEqual(this.service.expirationCalcTS, null, 'keeps expirationCalcTS as null');
    });
  });

  module('#mfaValidate', function () {
    test('it builds the mfa payload and calls the validate endpoint', async function (assert) {
      const mfaValidate = sinon.stub(this.service.api.sys, 'mfaValidate').resolves({ auth: {} });
      await this.service.mfaValidate({
        mfa_request_id: 'request-id',
        mfa_constraints: [
          { selectedMethod: { id: 'totp-id', type: 'totp' }, passcode: '123456' },
          { selectedMethod: { id: 'duo-id', type: 'duo' }, passcode: 'abc' },
          { selectedMethod: { id: 'duo-prefixed-id', type: 'duo' }, passcode: 'passcode=def' },
          { selectedMethod: { id: 'push-id', type: 'okta' }, passcode: '' },
        ],
      });
      assert.true(
        mfaValidate.calledOnceWithExactly({
          mfa_request_id: 'request-id',
          mfa_payload: {
            'totp-id': ['123456'],
            'duo-id': ['passcode=abc'],
            'duo-prefixed-id': ['passcode=def'],
            'push-id': [],
          },
        }),
        'passcode is only prefixed for duo and methods without a passcode send an empty payload'
      );
    });
  });

  module('#lookupTokenData', function () {
    test('it looks up the token with its own token header when data is missing', async function (assert) {
      const lookup = sinon
        .stub(this.service.api.auth, 'tokenLookUpSelf')
        .resolves({ data: { display_name: 'name', namespace_path: 'ns/', expire_time: 'time' } });
      const result = await this.service.lookupTokenData('my-token', true, {});
      assert.deepEqual(
        lookup.firstCall.args[0],
        this.service.api.buildHeaders({ token: 'my-token' }),
        'token is sent as the X-Vault-Token header'
      );
      assert.deepEqual(result, { displayName: 'name', namespacePath: 'ns/', expireTime: 'time' });
    });

    test('it skips the lookup when the login response already has the data', async function (assert) {
      const lookup = sinon.stub(this.service.api.auth, 'tokenLookUpSelf');
      const data = { displayName: 'name', namespacePath: 'ns/', expireTime: 'time' };
      assert.deepEqual(await this.service.lookupTokenData('my-token', true, data), data);
      assert.true(lookup.notCalled, 'lookup-self is not requested');
    });

    test('it falls back to the login data if the lookup fails', async function (assert) {
      sinon.stub(this.service.api.auth, 'tokenLookUpSelf').rejects(new Error('permission denied'));
      const data = { displayName: undefined, namespacePath: undefined, expireTime: undefined };
      assert.deepEqual(await this.service.lookupTokenData('my-token', true, data), data);
    });
  });
});
