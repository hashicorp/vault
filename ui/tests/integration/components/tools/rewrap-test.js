/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import { click, fillIn, find, render, waitUntil } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { TOOLS_SELECTORS as TS } from 'vault/tests/helpers/tools-selectors';
import sinon from 'sinon';

module('Integration | Component | tools/rewrap', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.renderComponent = async () => {
      await render(hbs`<Tools::Rewrap />`);
    };
  });

  test('it renders defaults', async function (assert) {
    await this.renderComponent();

    assert.dom('h1').hasText('Rewrap Token', 'title renders');
    assert.dom('label').hasText('Wrapped token');
    assert.dom(TS.toolsInput('original-token')).hasValue('');
    assert.dom(TS.toolsInput('rewrapped-token')).doesNotExist();
    assert.dom(TS.button('Done')).doesNotExist();
  });

  test('it renders errors', async function (assert) {
    this.server.post('sys/wrapping/rewrap', () => new Response(500, {}, { errors: ['Something is wrong'] }));
    await this.renderComponent();
    await click(TS.submit);
    await waitUntil(() => find(GENERAL.messageError));
    assert.dom(GENERAL.messageError).hasText('Error Something is wrong', 'Error renders');
  });

  test('it submits', async function (assert) {
    const flashSpy = sinon.spy(this.owner.lookup('service:flash-messages'), 'success');
    const original_token = 'original.OMZFbUurY0ppT2RTMGpRa0JOSUFqUzJUaGNqdWUQ6ooG=';
    const rewrapped_token = 'rewrapped_token.OMZFbUurY0ppT2RTMGpRa0JOSUFqUzJUaGNqdWUQ6ooG=';
    const data = { token: original_token };

    this.server.post('sys/wrapping/rewrap', (schema, req) => {
      assert.propEqual(JSON.parse(req.requestBody), data, `payload contains defaults: ${req.requestBody}`);
      return {
        wrap_info: {
          token: rewrapped_token,
          accessor: 'kfQad1FTIpXtdhWQMgpzcMFm',
          ttl: 1800,
          creation_time: '2024-06-05T13:57:28.827283-07:00',
          creation_path: 'sys/wrapping/wrap',
        },
      };
    });

    await this.renderComponent();

    // test submit
    await fillIn(TS.toolsInput('original-token'), original_token);
    await click(TS.submit);

    // test rewrapped token view
    await waitUntil(() => TS.toolsInput('rewrapped-token'));
    assert.true(flashSpy.calledWith('Rewrap was successful.'), 'it renders success flash');
    assert.dom('label').hasText('Rewrapped token');
    assert.dom(TS.toolsInput('rewrapped-token')).hasText(rewrapped_token);
    assert.dom(TS.toolsInput('original-token')).doesNotExist();

    // form resets clicking 'Done'
    await click(TS.button('Done'));
    assert.dom('label').hasText('Wrapped token');
    assert.dom(TS.toolsInput('original-token')).hasValue('', 'token input resets');
  });

  test('it trims token whitespace', async function (assert) {
    const data = { token: 'token.OMZFbUurY0ppT2RTMGpRa0JOSUFqUzJUaGNqdWUQ6ooG' };
    this.server.post('sys/wrapping/rewrap', (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.propEqual(payload, data, `token does not include whitespace: "${req.requestBody}"`);
      return {};
    });

    await this.renderComponent();

    await fillIn(TS.toolsInput('original-token'), `${data.token}  `);
    await click(TS.submit);
  });
});
