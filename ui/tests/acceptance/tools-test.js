/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, find, findAll, currentURL, visit, waitUntil } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { toolsActions } from 'vault/helpers/tools-actions';
import authPage from 'vault/tests/pages/auth';
import { capitalize } from '@ember/string';
import codemirror from 'vault/tests/helpers/codemirror';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { TOOLS_SELECTORS as TS } from 'vault/tests/helpers/tools-selectors';

const createTokenStore = () => {
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
const DATA_TO_WRAP = JSON.stringify({ tools: 'tests' });

module('Acceptance | tools', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    await authPage.login();
    return visit('/vault/tools');
  });

  test('it navigates to each action link', async function (assert) {
    assert.strictEqual(currentURL(), '/vault/tools/wrap', 'forwards from "vault/tools" to the first action');
    for (const action of toolsActions()) {
      await click(GENERAL.navLink(capitalize(action)));
      assert.strictEqual(currentURL(), `/vault/tools/${action}`, `it navigates to ${action}`);
    }
  });

  module('cross tool workflow', function () {
    test('it wraps data, performs lookup, rewraps and then unwraps data', async function (assert) {
      const tokenStore = createTokenStore();

      await waitUntil(() => find('.CodeMirror'));
      codemirror().setValue(DATA_TO_WRAP);

      await click(TS.submit);
      const wrappedToken = await waitUntil(() => find(TS.toolsInput('wrapping-token')));
      tokenStore.set(wrappedToken.innerText);

      // lookup
      await click(GENERAL.navLink('Lookup'));

      await fillIn(TS.toolsInput('wrapping-token'), tokenStore.get());
      await click(TS.submit);
      await waitUntil(() => findAll('[data-test-component="info-table-row"]').length >= 3);
      assert
        .dom(GENERAL.infoRowValue('Creation path'))
        .hasText('sys/wrapping/wrap', 'show creation path row');
      assert.dom(GENERAL.infoRowValue('Creation time')).exists();
      assert.dom(GENERAL.infoRowValue('Creation TTL')).hasText('1800', 'show creation ttl row');

      // rewrap
      await click(GENERAL.navLink('Rewrap'));

      await fillIn(TS.toolsInput('original-token'), tokenStore.get());
      await click(TS.submit);
      await waitUntil(() => find(TS.toolsInput('rewrapped-token')));
      const rewrappedToken = find(TS.toolsInput('rewrapped-token')).innerText;
      assert.notEqual(rewrappedToken, tokenStore.get(), 're-wrapped token is not the wrapped token');
      tokenStore.set(rewrappedToken);

      // unwrap
      await click(GENERAL.navLink('Unwrap'));

      await fillIn(TS.toolsInput('unwrap-token'), tokenStore.get());
      await click(TS.submit);
      await waitUntil(() => find('.CodeMirror'));
      assert.deepEqual(
        JSON.parse(codemirror().getValue()),
        JSON.parse(DATA_TO_WRAP),
        'unwrapped data equals input data'
      );
      await waitUntil(() => find(TS.tab('details')));
      await click(TS.tab('details'));
      await click(TS.tab('data'));
      assert.deepEqual(
        JSON.parse(codemirror().getValue()),
        JSON.parse(DATA_TO_WRAP),
        'data tab still has unwrapped data'
      );
    });
  });

  module('random', function () {
    test('it generates random bytes', async function (assert) {
      await click(GENERAL.navLink('Random'));
      assert.dom(TS.toolsInput('bytes')).hasValue('32', 'defaults to 32 bytes');
      await click(TS.submit);
      const randomBytes = await waitUntil(() => find(TS.toolsInput('random-bytes')));
      assert.strictEqual(randomBytes.innerText.length, 44, 'shows the returned value of random bytes');
    });
  });

  module('hash', function () {
    test('it generates hash', async function (assert) {
      await click(GENERAL.navLink('Hash'));

      await fillIn(TS.toolsInput('hash-input'), 'foo');
      await click(TS.toolsInput('b64-toggle'));
      assert.dom(TS.toolsInput('hash-input')).hasValue('Zm9v', 'it base64 encodes input');
      await click(TS.submit);
      let sumInput = await waitUntil(() => find(TS.toolsInput('sum')));
      assert
        .dom(sumInput)
        .hasText('LCa0a2j/xo/5m0U8HTBBNBNCLXBkg7+g+YpeiGJm564=', 'hashes the data, encodes input');
      await click(TS.button('Done'));

      await waitUntil(() => find(TS.toolsInput('hash-input')));
      assert.dom(TS.toolsInput('hash-input')).hasText('', 'it clears input on done');
      await fillIn(TS.toolsInput('hash-input'), 'e2RhdGE6ImZvbyJ9');

      await click(TS.submit);
      sumInput = await waitUntil(() => find(TS.toolsInput('sum')));
      assert
        .dom(sumInput)
        .hasText('JmSi2Hhbgu2WYOrcOyTqqMdym7KT3sohCwAwaMonVrc=', 'hashes the data, passes b64 input through');
    });
  });

  module('unwrap', function () {
    test('it unwraps with auth block', async function (assert) {
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
      this.server.post('/sys/wrapping/unwrap', () => {
        return AUTH_RESPONSE;
      });

      // unwrap
      await click(GENERAL.navLink('Unwrap'));

      await fillIn(TS.toolsInput('unwrap-token'), 'sometoken');
      await click(TS.submit);

      await waitUntil(() => find('.CodeMirror'));
      assert.deepEqual(
        JSON.parse(codemirror().getValue()),
        AUTH_RESPONSE.auth,
        'unwrapped data equals input data'
      );
    });
  });

  module('wrap', function () {
    test('it wraps data again after clicking "Back"', async function (assert) {
      const tokenStore = createTokenStore();
      await visit('/vault/tools/wrap');

      await waitUntil(() => find('.CodeMirror'));
      codemirror().setValue(DATA_TO_WRAP);

      // initial wrap
      await click(TS.submit);
      await waitUntil(() => find(TS.toolsInput('wrapping-token')));
      await click(TS.button('Back'));

      // wrap again without re-inputting data
      await click(TS.submit);
      const wrappedToken = await waitUntil(() => find(TS.toolsInput('wrapping-token')));
      tokenStore.set(wrappedToken.innerText);

      // there was a bug where clicking "back" cleared the parent's data, but not the child form component
      // so when users attempted to wrap data again the payload was actually empty and unwrapping the token returned {token: ""}
      // it is user desired behavior that the form does not clear on back, and that wrapping can be immediately repeated
      // we use lookup to check our token from the second wrap returns the unwrapped data we expect
      await click(GENERAL.navLink('Unwrap'));
      await fillIn(TS.toolsInput('unwrap-token'), tokenStore.get());
      await click(TS.submit);
      await waitUntil(() => find('.CodeMirror'));
      assert.strictEqual(codemirror().getValue(' '), '{   "tools": "tests" }', 'it renders unwrapped data');
    });

    test('it sends wrap ttl', async function (assert) {
      const tokenStore = createTokenStore();
      await visit('/vault/tools/wrap');

      await waitUntil(() => find('.CodeMirror'));
      codemirror().setValue(DATA_TO_WRAP);

      // update to non-default ttl
      await click(GENERAL.toggleInput('Wrap TTL'));
      await fillIn(GENERAL.ttl.input('Wrap TTL'), '20');

      await click(TS.submit);
      const wrappedToken = await waitUntil(() => find(TS.toolsInput('wrapping-token')));
      tokenStore.set(wrappedToken.innerText);

      // lookup to check ttl is what we expect
      await click(GENERAL.navLink('Lookup'));
      await fillIn(TS.toolsInput('wrapping-token'), tokenStore.get());
      await click(TS.submit);
      await waitUntil(() => findAll('[data-test-component="info-table-row"]').length >= 3);
      assert.dom(GENERAL.infoRowValue('Creation TTL')).hasText('1200', 'show creation ttl row');
    });
  });
});
