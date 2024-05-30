/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import {
  click,
  fillIn,
  find,
  findAll,
  currentURL,
  visit,
  settled,
  waitUntil,
  waitFor,
} from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { toolsActions } from 'vault/helpers/tools-actions';
import authPage from 'vault/tests/pages/auth';
import { capitalize } from '@ember/string';
import codemirror from 'vault/tests/helpers/codemirror';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { GENERAL } from '../helpers/general-selectors';

const SELECTORS = {
  submit: '[data-test-tools-submit]',
  toolsInput: (attr) => `[data-test-tools-input="${attr}"]`,
  tab: (item) => `[data-test-tab="${item}"]`,
  button: (action) => `[data-test-button="${action}"]`,
};
module('Acceptance | tools', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  const DATA_TO_WRAP = JSON.stringify({ tools: 'tests' });
  const TOOLS_ACTIONS = toolsActions();

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
      assert.dom(GENERAL.navLink(capitalize(action))).exists(`${action} link renders`);
    });

    await waitFor('.CodeMirror');
    codemirror().setValue(DATA_TO_WRAP);

    // wrap
    await click(SELECTORS.submit);
    const wrappedToken = await waitUntil(() => find(SELECTORS.toolsInput('wrapping-token')));
    tokenStore.set(wrappedToken.value);
    assert.dom(SELECTORS.toolsInput('wrapping-token')).hasValue(wrappedToken.value, 'has a wrapping token');

    //lookup
    await click(GENERAL.navLink('Lookup'));

    await fillIn(SELECTORS.toolsInput('wrapping-token'), tokenStore.get());
    await click(SELECTORS.submit);
    await waitUntil(() => findAll('[data-test-component="info-table-row"]').length >= 3);
    assert.dom(GENERAL.infoRowValue('Creation path')).hasText('sys/wrapping/wrap', 'show creation path row');
    assert.dom(GENERAL.infoRowValue('Creation time')).exists();
    assert.dom(GENERAL.infoRowValue('Creation TTL')).hasText('1800', 'show creation ttl row');

    //rewrap
    await click(GENERAL.navLink('Rewrap'));

    await fillIn(SELECTORS.toolsInput('wrapping-token'), tokenStore.get());
    await click(SELECTORS.submit);
    const rewrappedToken = await waitUntil(() => find(SELECTORS.toolsInput('rewrapped-token')));
    assert.ok(rewrappedToken.value, 'has a new re-wrapped token');
    assert.notEqual(rewrappedToken.value, tokenStore.get(), 're-wrapped token is not the wrapped token');
    tokenStore.set(rewrappedToken.value);
    await settled();

    //unwrap
    await click(GENERAL.navLink('Unwrap'));

    await fillIn(SELECTORS.toolsInput('wrapping-token'), tokenStore.get());
    await click(SELECTORS.submit);
    await waitFor('.CodeMirror');
    assert.deepEqual(
      JSON.parse(codemirror().getValue()),
      JSON.parse(DATA_TO_WRAP),
      'unwrapped data equals input data'
    );
    await waitUntil(() => find(SELECTORS.tab('details')));
    await click(SELECTORS.tab('details'));
    await click(SELECTORS.tab('data'));
    assert.deepEqual(
      JSON.parse(codemirror().getValue()),
      JSON.parse(DATA_TO_WRAP),
      'data tab still has unwrapped data'
    );
    //random
    await click(GENERAL.navLink('Random'));

    assert.dom(SELECTORS.toolsInput('bytes')).hasValue('32', 'defaults to 32 bytes');
    await click(SELECTORS.submit);
    const randomBytes = await waitUntil(() => find(SELECTORS.toolsInput('random-bytes')));
    assert.ok(randomBytes.value, 'shows the returned value of random bytes');

    //hash
    await click(GENERAL.navLink('Hash'));

    await fillIn(SELECTORS.toolsInput('hash-input'), 'foo');
    await click('[data-test-transit-b64-toggle="input"]');

    await click(SELECTORS.submit);
    let sumInput = await waitUntil(() => find(SELECTORS.toolsInput('sum')));
    assert
      .dom(sumInput)
      .hasValue('LCa0a2j/xo/5m0U8HTBBNBNCLXBkg7+g+YpeiGJm564=', 'hashes the data, encodes input');
    await click(SELECTORS.button('Back'));

    await fillIn(SELECTORS.toolsInput('hash-input'), 'e2RhdGE6ImZvbyJ9');

    await click(SELECTORS.submit);
    sumInput = await waitUntil(() => find(SELECTORS.toolsInput('sum')));
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
    this.server.post('/sys/wrapping/unwrap', () => {
      return AUTH_RESPONSE;
    });
    await visit('/vault/tools');

    //unwrap
    await click(GENERAL.navLink('Unwrap'));

    await fillIn(SELECTORS.toolsInput('wrapping-token'), 'sometoken');
    await click(SELECTORS.submit);

    await waitFor('.CodeMirror');
    assert.deepEqual(
      JSON.parse(codemirror().getValue()),
      AUTH_RESPONSE.auth,
      'unwrapped data equals input data'
    );
  });
});
