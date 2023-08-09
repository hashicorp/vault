/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { click, fillIn, find, findAll, currentURL, visit, settled, waitUntil } from '@ember/test-helpers';
import Pretender from 'pretender';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { toolsActions } from 'vault/helpers/tools-actions';
import authPage from 'vault/tests/pages/auth';
import { capitalize } from '@ember/string';

module('Acceptance | tools', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  const DATA_TO_WRAP = JSON.stringify({ tools: 'tests' });
  const TOOLS_ACTIONS = toolsActions();

  /*
  data-test-tools-input="wrapping-token"
  data-test-tools-input="rewrapped-token"
  data-test-tools="token-lookup-row"
  data-test-sidebar-nav-link=supportedAction
  */

  var createTokenStore = () => {
    let token;
    return {
      set(val) {
        token = val;
      },
      get() {
        return token;
      },
    };
  };
  test('tools functionality', async function (assert) {
    var tokenStore = createTokenStore();
    await visit('/vault/tools');

    assert.strictEqual(currentURL(), '/vault/tools/wrap', 'forwards to the first action');
    TOOLS_ACTIONS.forEach((action) => {
      assert.dom(`[data-test-sidebar-nav-link="${capitalize(action)}"]`).exists(`${action} link renders`);
    });

    const { CodeMirror } = await waitUntil(() => find('.CodeMirror'));
    CodeMirror.setValue(DATA_TO_WRAP);

    // wrap
    await click('[data-test-tools-submit]');
    const wrappedToken = await waitUntil(() => find('[data-test-tools-input="wrapping-token"]'));
    tokenStore.set(wrappedToken.value);
    assert
      .dom('[data-test-tools-input="wrapping-token"]')
      .hasValue(wrappedToken.value, 'has a wrapping token');

    //lookup
    await click('[data-test-sidebar-nav-link="Lookup"]');

    await fillIn('[data-test-tools-input="wrapping-token"]', tokenStore.get());
    await click('[data-test-tools-submit]');
    await waitUntil(() => findAll('[data-test-tools="token-lookup-row"]').length >= 3);
    const rows = findAll('[data-test-tools="token-lookup-row"]');
    assert.dom(rows[0]).hasText(/Creation path/, 'show creation path row');
    assert.dom(rows[1]).hasText(/Creation time/, 'show creation time row');
    assert.dom(rows[2]).hasText(/Creation TTL/, 'show creation ttl row');

    //rewrap
    await click('[data-test-sidebar-nav-link="Rewrap"]');

    await fillIn('[data-test-tools-input="wrapping-token"]', tokenStore.get());
    await click('[data-test-tools-submit]');
    const rewrappedToken = await waitUntil(() => find('[data-test-tools-input="rewrapped-token"]'));
    assert.ok(rewrappedToken.value, 'has a new re-wrapped token');
    assert.notEqual(rewrappedToken.value, tokenStore.get(), 're-wrapped token is not the wrapped token');
    tokenStore.set(rewrappedToken.value);
    await settled();

    //unwrap
    await click('[data-test-sidebar-nav-link="Unwrap"]');

    await fillIn('[data-test-tools-input="wrapping-token"]', tokenStore.get());
    await click('[data-test-tools-submit]');
    assert.deepEqual(
      JSON.parse(CodeMirror.getValue()),
      JSON.parse(DATA_TO_WRAP),
      'unwrapped data equals input data'
    );
    const buttonDetails = await waitUntil(() => find('[data-test-button-details]'));
    await click(buttonDetails);
    await click('[data-test-button-data]');
    assert.dom('.CodeMirror').exists();

    //random
    await click('[data-test-sidebar-nav-link="Random"]');

    assert.dom('[data-test-tools-input="bytes"]').hasValue('32', 'defaults to 32 bytes');
    await click('[data-test-tools-submit]');
    const randomBytes = await waitUntil(() => find('[data-test-tools-input="random-bytes"]'));
    assert.ok(randomBytes.value, 'shows the returned value of random bytes');

    //hash
    await click('[data-test-sidebar-nav-link="Hash"]');

    await fillIn('[data-test-tools-input="hash-input"]', 'foo');
    await click('[data-test-transit-b64-toggle="input"]');

    await click('[data-test-tools-submit]');
    let sumInput = await waitUntil(() => find('[data-test-tools-input="sum"]'));
    assert
      .dom(sumInput)
      .hasValue('LCa0a2j/xo/5m0U8HTBBNBNCLXBkg7+g+YpeiGJm564=', 'hashes the data, encodes input');
    await click('[data-test-tools-back]');

    await fillIn('[data-test-tools-input="hash-input"]', 'e2RhdGE6ImZvbyJ9');

    await click('[data-test-tools-submit]');
    sumInput = await waitUntil(() => find('[data-test-tools-input="sum"]'));
    assert
      .dom(sumInput)
      .hasValue('JmSi2Hhbgu2WYOrcOyTqqMdym7KT3sohCwAwaMonVrc=', 'hashes the data, passes b64 input through');
  });

  const AUTH_RESPONSE = {
    request_id: '39802bc4-235c-2f0b-87f3-ccf38503ac3e',
    lease_id: '',
    renewable: false,
    lease_duration: 0,
    data: null,
    wrap_info: null,
    warnings: null,
    auth: {
      client_token: 'ecfc2758-588e-981d-50f4-a25883bbf03c',
      accessor: '6299780b-f2b2-1a3f-7b83-9d3d67629249',
      policies: ['root'],
      metadata: null,
      lease_duration: 0,
      renewable: false,
      entity_id: '',
    },
  };

  test('ensure unwrap with auth block works properly', async function (assert) {
    this.server = new Pretender(function () {
      this.post('/v1/sys/wrapping/unwrap', (response) => {
        return [response, { 'Content-Type': 'application/json' }, JSON.stringify(AUTH_RESPONSE)];
      });
    });
    await visit('/vault/tools');

    //unwrap
    await click('[data-test-sidebar-nav-link="Unwrap"]');

    await fillIn('[data-test-tools-input="wrapping-token"]', 'sometoken');
    await click('[data-test-tools-submit]');

    assert.deepEqual(
      JSON.parse(findAll('.CodeMirror')[0].CodeMirror.getValue()),
      AUTH_RESPONSE.auth,
      'unwrapped data equals input data'
    );
    this.server.shutdown();
  });
});
