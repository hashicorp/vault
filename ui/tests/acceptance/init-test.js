import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';

import initPage from 'vault/tests/pages/init';
import Pretender from 'pretender';

const HEALTH_RESPONSE = {
  initialized: false,
  sealed: true,
  standby: true,
  performance_standby: false,
  replication_performance_mode: 'unknown',
  replication_dr_mode: 'unknown',
  server_time_utc: 1538066726,
  version: '0.11.0+prem',
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

module('Acceptance | init', function(hooks) {
  setupApplicationTest(hooks);

  let setInitResponse = (server, resp) => {
    server.put('/v1/sys/init', () => {
      return [200, { 'Content-Type': 'application/json' }, JSON.stringify(resp)];
    });
  };
  let setStatusResponse = (server, resp) => {
    server.get('/v1/sys/seal-status', () => {
      return [200, { 'Content-Type': 'application/json' }, JSON.stringify(resp)];
    });
  };
  hooks.beforeEach(function() {
    this.server = new Pretender();
    this.server.get('/v1/sys/health', () => {
      return [200, { 'Content-Type': 'application/json' }, JSON.stringify(HEALTH_RESPONSE)];
    });
  });

  hooks.afterEach(function() {
    this.server.shutdown();
  });

  test('cloud seal init', async function(assert) {
    setInitResponse(this.server, CLOUD_SEAL_RESPONSE);
    setStatusResponse(this.server, CLOUD_SEAL_STATUS_RESPONSE);
    await initPage.init(5, 3);
    assert.equal(
      initPage.keys.length,
      CLOUD_SEAL_RESPONSE.recovery_keys.length,
      'shows all of the recovery keys'
    );
    assert.equal(initPage.buttonText, 'Continue to Authenticate', 'links to authenticate');
    let { requestBody } = this.server.handledRequests.findBy('url', '/v1/sys/init');
    requestBody = JSON.parse(requestBody);
    for (let attr of ['recovery_shares', 'recovery_threshold']) {
      assert.ok(requestBody[attr], `requestBody includes cloud seal specific attribute: ${attr}`);
    }
  });

  test('shamir seal init', async function(assert) {
    setInitResponse(this.server, SEAL_RESPONSE);
    setStatusResponse(this.server, SEAL_STATUS_RESPONSE);

    await initPage.init(3, 2);
    assert.equal(initPage.keys.length, SEAL_RESPONSE.keys.length, 'shows all of the recovery keys');
    assert.equal(initPage.buttonText, 'Continue to Unseal', 'links to unseal');

    let { requestBody } = this.server.handledRequests.findBy('url', '/v1/sys/init');
    requestBody = JSON.parse(requestBody);
    for (let attr of ['recovery_shares', 'recovery_threshold']) {
      assert.notOk(requestBody[attr], `requestBody does not include cloud seal specific attribute: ${attr}`);
    }
  });
});
