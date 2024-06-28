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

module('Integration | Component | tools/random', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.random_bytes = null;
    this.bytes = 32;
    this.format = 'base64';
    this.renderComponent = async () => {
      await render(hbs`<Tools::Random />`);
    };
  });

  test('it renders defaults', async function (assert) {
    await this.renderComponent();

    assert.dom('h1').hasText('Random Bytes', 'Title renders');
    assert.dom('#bytes').hasValue('32');
    assert.dom('#format').hasValue('base64');
    assert.dom(TS.submit).hasText('Generate');
    assert.dom(TS.toolsInput('random-bytes')).doesNotExist();
    assert.dom(TS.button('Done')).doesNotExist();
  });

  test('it renders errors', async function (assert) {
    this.server.post('sys/tools/random', () => new Response(500, {}, { errors: ['Something is wrong'] }));
    await this.renderComponent();
    await click(TS.submit);
    await waitUntil(() => find(GENERAL.messageError));
    assert.dom(GENERAL.messageError).hasText('Error Something is wrong');
  });

  test('it submits with default values', async function (assert) {
    const flashSpy = sinon.spy(this.owner.lookup('service:flash-messages'), 'success');
    const random_bytes = 'ekZtsOvSOvKEPUmJohRj1nwONQ1XafwcXGZEwy/0nbY==';
    const data = { format: 'base64', bytes: 32 };

    this.server.post('sys/tools/random', (schema, req) => {
      assert.propEqual(JSON.parse(req.requestBody), data, `payload contains defaults: ${req.requestBody}`);
      return { data: { random_bytes } };
    });

    await this.renderComponent();

    // test submit
    await click(TS.submit);

    // test random bytes view
    await waitUntil(() => TS.toolsInput('random-bytes'));
    assert.true(flashSpy.calledWith('Generated random bytes successfully.'), 'it renders success flash');
    assert.dom(TS.toolsInput('random-bytes')).hasText(random_bytes);
    assert.dom('label').hasText('Random bytes');
    assert.dom('#bytes').doesNotExist();
    assert.dom('#format').doesNotExist();

    // clicking 'Done' resets form
    await click(TS.button('Done'));
    assert.dom('#bytes').hasValue('32');
    assert.dom('#format').hasValue('base64');
  });

  test('it submits with updated values', async function (assert) {
    const random_bytes =
      '23d0cfaf42c93afd878e77ee484bf0605ce5294459ef755f5780e77e5e162a27e64e3ebda62152ac88035d64';
    const data = { format: 'hex', bytes: 44 };

    this.server.post('sys/tools/random', (schema, req) => {
      assert.propEqual(JSON.parse(req.requestBody), data, `payload has updated data: ${req.requestBody}`);
      return { data: { random_bytes } };
    });

    await this.renderComponent();
    await fillIn('#bytes', '44');
    await fillIn('#format', 'hex');
    await click(TS.submit);
  });
});
