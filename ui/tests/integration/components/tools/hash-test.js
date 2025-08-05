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

module('Integration | Component | tools/hash', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.renderComponent = async () => {
      await render(hbs`<Tools::Hash />`);
    };
  });

  test('it renders form', async function (assert) {
    await this.renderComponent();

    assert.dom('h1').hasText('Hash Data');
    assert.dom(TS.toolsInput('hash-input')).hasValue('');
    assert.dom('#algorithm').hasValue('sha2-256');
    assert.dom('#format').hasValue('base64');
    assert.dom(TS.toolsInput('sum')).doesNotExist();
    assert.dom(TS.button('Done')).doesNotExist();
  });

  test('it renders errors', async function (assert) {
    this.server.post('sys/tools/hash', () => new Response(500, {}, { errors: ['Something is wrong'] }));
    await this.renderComponent();
    await click(TS.submit);
    await waitUntil(() => find(GENERAL.messageError));
    assert.dom(GENERAL.messageError).hasText('Error Something is wrong');
  });

  test('it submits with default values', async function (assert) {
    const flashSpy = sinon.spy(this.owner.lookup('service:flash-messages'), 'success');
    const sum = 'GmoKPULUXifIFSGPZx29CGSm8MFnBuk4SGPsmFlduGc=';
    const data = {
      algorithm: 'sha2-256',
      format: 'base64',
      input: 'blah',
    };

    this.server.post('sys/tools/hash', (schema, req) => {
      assert.propEqual(JSON.parse(req.requestBody), data, `payload contains defaults: ${req.requestBody}`);
      return { data: { sum } };
    });

    await this.renderComponent();

    // test submit
    await fillIn(TS.toolsInput('hash-input'), 'blah');
    await click(TS.submit);

    // test sum view
    await waitUntil(() => TS.toolsInput('sum'));
    assert.true(flashSpy.calledWith('Hash was successful.'), 'it renders success flash');
    assert.dom(TS.toolsInput('sum')).hasText(sum);
    assert.dom('label').hasText('Sum');
    assert.dom('#algorithm').doesNotExist();
    assert.dom('#format').doesNotExist();

    // test form reset clicking 'Done'
    await click(TS.button('Done'));
    assert.dom('#algorithm').hasValue('sha2-256');
    assert.dom('#format').hasValue('base64');
    assert.dom(TS.toolsInput('hash-input')).hasValue('', 'inputs reset to default values');
  });

  test('it submits with updated values', async function (assert) {
    const sum = '07a49af6947eaa5ddce0d40aa4584687a95f9c0be0d1d7df009d63da=';
    const data = {
      algorithm: 'sha2-224',
      format: 'hex',
      input: 'blah',
    };

    this.server.post('sys/tools/hash', (schema, req) => {
      assert.propEqual(JSON.parse(req.requestBody), data, `payload has updated data: ${req.requestBody}`);
      return { data: { sum } };
    });

    await this.renderComponent();

    // submits updated values
    await fillIn(TS.toolsInput('hash-input'), 'blah');
    await fillIn('#algorithm', 'sha2-224');
    await fillIn('#format', 'hex');
    await click(TS.submit);
  });
});
