/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

module('Integration | Component | raft-storage-overview', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.model = [
      { address: '127.0.0.1:8200', voter: true },
      { address: '127.0.0.1:8200', voter: true, leader: true },
    ];
  });

  test('it renders', async function (assert) {
    await render(hbs`<RaftStorageOverview @model={{this.model}} />`);
    assert.dom('[data-raft-row]').exists({ count: 2 });
  });

  test('it should download snapshot via service worker', async function (assert) {
    assert.expect(3);

    const token = this.owner.lookup('service:auth').currentToken;
    const generateMockEvent = (action) => ({
      data: { action },
      ports: [
        {
          postMessage(message) {
            const getToken = action === 'getToken';
            const expected = getToken ? { token } : { error: 'Unknown request' };
            assert.deepEqual(
              message,
              expected,
              `${
                getToken ? 'Token' : 'Error'
              } is returned to service worker in message event listener callback`
            );
          },
        },
      ],
    });

    sinon.stub(navigator.serviceWorker, 'getRegistration').resolves(true);
    sinon.stub(navigator.serviceWorker, 'addEventListener').callsFake((name, cb) => {
      assert.strictEqual(name, 'message', 'Event listener added for service worker message');
      cb(generateMockEvent('getToken'));
      cb(generateMockEvent('unknown'));
    });

    await render(hbs`<RaftStorageOverview @model={{this.model}} />`);
    // avoid clicking the download button or the url will change
    // the service worker invokes the event listener callback when it intercepts the request to /v1/sys/storage/raft/snapshot
    // for the test we manually fire the callback as soon as it is passed to the addEventListener stub
  });
});
