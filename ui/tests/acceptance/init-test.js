/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';

import initPage from 'vault/tests/pages/init';
import Pretender from 'pretender';
import { setRunOptions } from 'ember-a11y-testing/test-support';

const HEALTH_RESPONSE = {
  initialized: false,
  sealed: true,
  standby: true,
  performance_standby: false,
  replication_performance_mode: 'unknown',
  replication_dr_mode: 'unknown',
  server_time_utc: 1538066726,
  version: '1.13.0-dev1',
};

const CLOUD_SEAL_RESPONSE = {
  keys: [],
  keys_base64: [],
  recovery_keys: [
    '1659986a8d56b998b175b6e259998f3c064c061d256c2a331681b8d122fedf0db4',
    '4d34c58f56e4f077e3b74f9e8db2850fc251ac3f16e952441301eedc462addeb84',
    '3b3cbdf4b2f5ac1e809ff1bb72fd9778e460856561728a871a9370345bd52e97f4',
    'aa99b46e2ed5d837ee9824b7894b24987be2f32c81ab9ff5ce9e07d2012eaf4158',
    'c2bf6d71d8db8ae09b26177ed393ecb274740fe9ab51884eaa00ac113a74c08ba7',
  ],
  recovery_keys_base64: [
    'FlmYao1WuZixdbbiWZmPPAZMBh0lbCozFoG40SL+3w20',
    'TTTFj1bk8Hfjt0+ejbKFD8JRrD8W6VJEEwHu3EYq3euE',
    'Ozy99LL1rB6An/G7cv2XeORghWVhcoqHGpNwNFvVLpf0',
    'qpm0bi7V2DfumCS3iUskmHvi8yyBq5/1zp4H0gEur0FY',
    'wr9tcdjbiuCbJhd+05PssnR0D+mrUYhOqgCsETp0wIun',
  ],
  root_token: '48dF3Drr1jl4ayM0jcHrN4NC',
};
const SEAL_RESPONSE = {
  keys: [
    '1659986a8d56b998b175b6e259998f3c064c061d256c2a331681b8d122fedf0db4',
    '4d34c58f56e4f077e3b74f9e8db2850fc251ac3f16e952441301eedc462addeb84',
    '3b3cbdf4b2f5ac1e809ff1bb72fd9778e460856561728a871a9370345bd52e97f4',
  ],
  keys_base64: [
    'FlmYao1WuZixdbbiWZmPPAZMBh0lbCozFoG40SL+3w20',
    'TTTFj1bk8Hfjt0+ejbKFD8JRrD8W6VJEEwHu3EYq3euE',
    'Ozy99LL1rB6An/G7cv2XeORghWVhcoqHGpNwNFvVLpf0',
  ],
  root_token: '48dF3Drr1jl4ayM0jcHrN4NC',
};

const CLOUD_SEAL_STATUS_RESPONSE = {
  type: 'awskms',
  sealed: true,
  initialized: false,
};
const SEAL_STATUS_RESPONSE = {
  type: 'shamir',
  sealed: true,
  initialized: false,
};

const assertRequest = (req, assert, isCloud) => {
  const json = JSON.parse(req.requestBody);
  for (const key of ['recovery_shares', 'recovery_threshold']) {
    assert[isCloud ? 'ok' : 'notOk'](
      json[key],
      `requestBody ${isCloud ? 'includes' : 'does not include'} cloud seal specific attribute: ${key}`
    );
  }
  for (const key of ['secret_shares', 'secret_threshold']) {
    assert[isCloud ? 'notOk' : 'ok'](
      json[key],
      `requestBody ${isCloud ? 'does not include' : 'includes'} shamir specific attribute: ${key}`
    );
  }
};

module('Acceptance | init', function (hooks) {
  setupApplicationTest(hooks);

  const setInitResponse = (server, resp) => {
    server.put('/v1/sys/init', () => {
      return [200, { 'Content-Type': 'application/json' }, JSON.stringify(resp)];
    });
  };
  const setStatusResponse = (server, resp) => {
    server.get('/v1/sys/seal-status', () => {
      return [200, { 'Content-Type': 'application/json' }, JSON.stringify(resp)];
    });
  };
  hooks.beforeEach(function () {
    this.server = new Pretender();
    this.server.get('/v1/sys/health', () => {
      return [200, { 'Content-Type': 'application/json' }, JSON.stringify(HEALTH_RESPONSE)];
    });
    this.server.get('/v1/sys/internal/ui/feature-flags', this.server.passthrough);
  });

  hooks.afterEach(function () {
    this.server.shutdown();
  });

  test('cloud seal init', async function (assert) {
    // continue button is disabled, violating color-contrast
    setRunOptions({
      rules: {
        'color-contrast': { enabled: false },
      },
    });
    assert.expect(6);

    setInitResponse(this.server, CLOUD_SEAL_RESPONSE);
    setStatusResponse(this.server, CLOUD_SEAL_STATUS_RESPONSE);

    await initPage.init(5, 3);

    assert.strictEqual(
      initPage.keys.length,
      CLOUD_SEAL_RESPONSE.recovery_keys.length,
      'shows all of the recovery keys'
    );
    assert.strictEqual(initPage.buttonText, 'Continue to Authenticate', 'links to authenticate');
    assertRequest(this.server.handledRequests.findBy('url', '/v1/sys/init'), assert, true);
  });

  test('shamir seal init', async function (assert) {
    assert.expect(6);

    setInitResponse(this.server, SEAL_RESPONSE);
    setStatusResponse(this.server, SEAL_STATUS_RESPONSE);

    await initPage.init(3, 2);

    assert.strictEqual(initPage.keys.length, SEAL_RESPONSE.keys.length, 'shows all of the recovery keys');
    assert.strictEqual(initPage.buttonText, 'Continue to Unseal', 'links to unseal');
    assertRequest(this.server.handledRequests.findBy('url', '/v1/sys/init'), assert, false);
  });
});
