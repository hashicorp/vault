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
import { TTL_PICKER as TTL } from 'vault/tests/helpers/components/ttl-picker-selectors';
import { TOOLS_SELECTORS as TS } from 'vault/tests/helpers/tools-selectors';
import codemirror from 'vault/tests/helpers/codemirror';
import sinon from 'sinon';

module('Integration | Component | tools/wrap', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.renderComponent = async () => {
      await render(hbs`
    <Tools::Wrap />`);
    };
    this.wrapData = `{"foo": "bar"}`;
    this.token = 'blah.jhfel7SmsVeZwihaGiIKHGh2cy5XZWtEeEt5WmRwS1VYSTNDb1BBVUNsVFAQ3JIK';
    // default mirage response here is overridden in some tests
    this.server.post('sys/wrapping/wrap', () => {
      // removed superfluous response data for this test
      return { wrap_info: { token: this.token } };
    });
  });

  test('it renders defaults', async function (assert) {
    await this.renderComponent();

    assert.dom('h1').hasText('Wrap Data', 'Title renders');
    assert.dom('label').hasText('Data to wrap (json-formatted)');
    assert.strictEqual(codemirror().getValue(' '), '{ }', 'json editor initializes with empty object');
    assert.dom(TTL.toggleByLabel('Wrap TTL')).isNotChecked('Wrap TTL defaults to unchecked');
    assert.dom(TS.submit).isEnabled();
    assert.dom(TS.toolsInput('wrapping-token')).doesNotExist();
    assert.dom(TS.button('Back')).doesNotExist();
    assert.dom(TS.button('Done')).doesNotExist();

    await click(TTL.toggleByLabel('Wrap TTL'));
    assert.dom(TTL.valueInputByLabel('Wrap TTL')).hasValue('30', 'ttl defaults to 30 when toggled');
    assert.dom(TTL.ttlUnit).hasValue('m', 'ttl defaults to minutes when toggled');
  });

  test('it renders errors', async function (assert) {
    this.server.post('sys/wrapping/wrap', () => new Response(500, {}, { errors: ['Something is wrong'] }));
    await this.renderComponent();
    await click(TS.submit);
    await waitUntil(() => find(GENERAL.messageError));
    assert.dom(GENERAL.messageError).hasText('Error Something is wrong', 'Error renders');
  });

  test('it submits with defaults', async function (assert) {
    assert.expect(6);
    const flashSpy = sinon.spy(this.owner.lookup('service:flash-messages'), 'success');

    this.server.post('sys/wrapping/wrap', (schema, { requestBody, requestHeaders }) => {
      const payload = JSON.parse(requestBody);
      assert.propEqual(payload, JSON.parse(this.wrapData), `payload contains data: ${requestBody}`);
      assert.strictEqual(requestHeaders['X-Vault-Wrap-TTL'], '30m', 'request header has default wrap ttl');
      return {
        wrap_info: {
          token: this.token,
          accessor: '5yjKx6Om9NmBx1mjiN1aIrnm',
          ttl: 1800,
          creation_time: '2024-06-07T12:02:22.096254-07:00',
          creation_path: 'sys/wrapping/wrap',
        },
      };
    });

    await this.renderComponent();
    await codemirror().setValue(this.wrapData);
    await click(TS.submit);
    await waitUntil(() => find(TS.toolsInput('wrapping-token')));
    assert.true(flashSpy.calledWith('Wrap was successful.'), 'it renders success flash');
    assert.dom(TS.toolsInput('wrapping-token')).hasText(this.token);
    assert.dom('label').hasText('Wrapped token');
    assert.dom('.CodeMirror').doesNotExist();
  });

  test('it submits with updated ttl', async function (assert) {
    assert.expect(2);
    this.server.post('sys/wrapping/wrap', (schema, { requestBody, requestHeaders }) => {
      const payload = JSON.parse(requestBody);
      assert.propEqual(payload, JSON.parse(this.wrapData), `payload contains data: ${requestBody}`);
      assert.strictEqual(requestHeaders['X-Vault-Wrap-TTL'], '1200s', 'request header has updated wrap ttl');
      // only testing payload/header assertions, no need for return here
      return {};
    });

    await this.renderComponent();
    await codemirror().setValue(this.wrapData);
    await click(TTL.toggleByLabel('Wrap TTL'));
    await fillIn(TTL.valueInputByLabel('Wrap TTL'), '20');
    await click(TS.submit);
  });

  test('it resets on done', async function (assert) {
    await this.renderComponent();
    await codemirror().setValue(this.wrapData);
    await click(TTL.toggleByLabel('Wrap TTL'));
    await fillIn(TTL.valueInputByLabel('Wrap TTL'), '20');
    await click(TS.submit);

    await waitUntil(() => find(TS.button('Done')));
    await click(TS.button('Done'));
    assert.strictEqual(codemirror().getValue(' '), '{ }', 'json editor resets to empty object');
    assert.dom(TTL.toggleByLabel('Wrap TTL')).isNotChecked('Wrap TTL resets to unchecked');
    await click(TTL.toggleByLabel('Wrap TTL'));
    assert.dom(TTL.valueInputByLabel('Wrap TTL')).hasValue('30', 'ttl resets to default when toggled');
  });

  test('it preserves input data on back', async function (assert) {
    await this.renderComponent();
    await codemirror().setValue(this.wrapData);
    await click(TS.submit);

    await waitUntil(() => find(TS.button('Back')));
    await click(TS.button('Back'));
    assert.strictEqual(codemirror().getValue(' '), `{"foo": "bar"}`, 'json editor has original data');
    assert.dom(TTL.toggleByLabel('Wrap TTL')).isNotChecked('Wrap TTL defaults to unchecked');
  });

  test('it disables/enables submit based on json linting', async function (assert) {
    await this.renderComponent();
    await codemirror().setValue(`{bad json}`);
    assert.dom(TS.submit).isDisabled('submit disables if json editor has linting errors');

    await codemirror().setValue(this.wrapData);
    assert.dom(TS.submit).isEnabled('submit reenables if json editor has no linting errors');
  });
});
