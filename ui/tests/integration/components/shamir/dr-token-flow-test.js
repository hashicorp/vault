/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import sinon from 'sinon';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, find, render, waitFor, waitUntil } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SHAMIR_FORM } from 'vault/tests/helpers/components/shamir-selectors';

module('Integration | Component | shamir/dr-token-flow', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  test('begin to middle flow works', async function (assert) {
    assert.expect(16);
    this.server.get('/sys/replication/dr/secondary/generate-operation-token/attempt', function () {
      assert.ok('Check endpoint is queried');
      return {};
    });
    this.server.post('/sys/replication/dr/secondary/generate-operation-token/attempt', function (_, req) {
      const requestBody = JSON.parse(req.requestBody);
      assert.ok('Starts the token generation');
      assert.deepEqual(requestBody, { attempt: true });
      return {
        started: true,
        nonce: 'nonce-1234',
        progress: 0,
        required: 3,
        encoded_token: '',
        otp: 'otp-9876',
        otp_length: 24,
        complete: false,
      };
    });
    this.server.post('/sys/replication/dr/secondary/generate-operation-token/update', function (_, req) {
      const requestBody = JSON.parse(req.requestBody);
      assert.ok('Makes request at the /update path');
      assert.deepEqual(requestBody, { key: 'some-key', nonce: 'nonce-1234' });
      return {
        started: true,
        nonce: 'nonce-1234',
        progress: 1,
        required: 3,
        encoded_token: '',
        otp: '',
        otp_length: 24,
        complete: false,
      };
    });

    await render(hbs`<Shamir::DrTokenFlow @action="generate-dr-operation-token" />`);
    assert.dom(SHAMIR_FORM.flowStep('begin')).exists('First step shows');
    assert.dom(GENERAL.button('use-pgp-key-cta')).hasText('Provide PGP Key');
    assert.dom(GENERAL.button('generate-token-cta')).hasText('Generate operation token');

    await click(GENERAL.button('generate-token-cta'));
    assert.ok(
      await waitUntil(() => find(SHAMIR_FORM.flowStep('primary-token'))),
      'shows primary token step after start'
    );

    await fillIn(GENERAL.inputByAttr('primary-token'), 'some-token');
    await click('[data-test-submit-primary-token]');
    assert.ok(
      await waitUntil(() => find(SHAMIR_FORM.flowStep('shamir'))),
      'shows shamir step after primary token input'
    );
    assert
      .dom(SHAMIR_FORM.progress)
      .hasText('0/3 keys provided', 'progress shows reflecting checkStatus response with defaults');
    assert.dom(SHAMIR_FORM.otpInfo).exists('OTP info banner shows');
    assert.dom(SHAMIR_FORM.otpCode).hasText('otp-9876', 'Shows OTP in copy banner');
    // Fill in shamir key and submit
    await fillIn(GENERAL.inputByAttr('shamir-key'), 'some-key');
    await click(GENERAL.submitButton);

    assert.ok(
      await waitUntil(() => find(SHAMIR_FORM.otpInfo)),
      'OTP info still banner shows even when attempt response does not include it'
    );
    assert
      .dom(SHAMIR_FORM.otpCode)
      .hasText('otp-9876', 'Still shows OTP in copy banner when attempt response does not include it');
    assert
      .dom(SHAMIR_FORM.progress)
      .hasText('1/3 keys provided', 'progress shows reflecting attempt response');
  });

  test('middle to finish flow works', async function (assert) {
    assert.expect(10);
    this.server.get('/sys/replication/dr/secondary/generate-operation-token/attempt', function () {
      assert.ok('Check endpoint is queried');
      return {
        started: true,
        nonce: 'nonce-1234',
        progress: 2,
        required: 3,
        encoded_token: '',
        otp: '',
        otp_length: 24,
        complete: false,
      };
    });
    this.server.post('/sys/replication/dr/secondary/generate-operation-token/attempt', function () {
      assert.notOk('attempt endpoint should not be queried');
    });
    this.server.post('/sys/replication/dr/secondary/generate-operation-token/update', function (_, req) {
      const requestBody = JSON.parse(req.requestBody);
      assert.ok('Makes request at the /update path');
      assert.deepEqual(requestBody, { key: 'some-key', nonce: 'nonce-1234' });
      return {
        started: true,
        nonce: 'nonce-1234',
        progress: 3,
        required: 3,
        encoded_token: 'encoded-token-here',
        otp: '',
        otp_length: 24,
        complete: true,
      };
    });
    await render(hbs`<Shamir::DrTokenFlow @action="generate-dr-operation-token" />`);
    await click(GENERAL.button('generate-token-cta'));
    assert.ok(
      await waitUntil(() => find(SHAMIR_FORM.flowStep('primary-token'))),
      'shows primary token step after start'
    );

    await fillIn(GENERAL.inputByAttr('primary-token'), 'some-token');
    await click('[data-test-submit-primary-token]');
    assert.ok(
      await waitUntil(() => find(SHAMIR_FORM.flowStep('shamir'))),
      'shows shamir step after primary token input'
    );
    assert
      .dom(SHAMIR_FORM.progress)
      .hasText('2/3 keys provided', 'progress shows reflecting checkStatus response');
    assert.dom(SHAMIR_FORM.otpInfo).doesNotExist('OTP info banner not shown');
    assert.dom(SHAMIR_FORM.otpCode).doesNotExist('otp-9876', 'OTP copy banner not shown');
    await fillIn(GENERAL.inputByAttr('shamir-key'), 'some-key');
    await click(GENERAL.submitButton);

    assert.ok(
      await waitUntil(() => find(SHAMIR_FORM.flowStep('show-token'))),
      'updates to show encoded token on complete'
    );
    assert
      .dom(GENERAL.copySnippet('shamir-encoded-token'))
      .hasText('encoded-token-here', 'shows encoded token from /update response');
  });

  test('it works correctly when pgp key chosen', async function (assert) {
    assert.expect(4);
    this.server.get('/sys/replication/dr/secondary/generate-operation-token/attempt', function () {
      return {};
    });
    this.server.post(
      '/sys/replication/dr/secondary/generate-operation-token/attempt',
      function (schema, req) {
        const body = JSON.parse(req.requestBody);
        assert.deepEqual(body, { pgp_key: 'some-key-here' }, 'correct payload');
        return {
          started: true,
          progress: 1,
          required: 3,
          complete: false,
        };
      }
    );
    await render(hbs`<Shamir::DrTokenFlow @action="generate-dr-operation-token" />`);
    await click(GENERAL.button('use-pgp-key-cta'));

    assert.ok(await waitUntil(() => find(SHAMIR_FORM.flowStep('choose-pgp-key'))), 'PGP form shows');
    await click(GENERAL.textToggle);
    await fillIn(GENERAL.textareaByAttr('pgp-key'), 'some-key-here');
    await click(GENERAL.button('use-pgp-key'));
    await click(GENERAL.submitButton);

    assert.ok(
      await waitUntil(() => find(SHAMIR_FORM.flowStep('primary-token'))),
      'shows primary token input after start'
    );
    await fillIn(GENERAL.inputByAttr('primary-token'), 'some-token');
    await click('[data-test-submit-primary-token]');

    assert.ok(
      await waitUntil(() => find(SHAMIR_FORM.flowStep('shamir'))),
      'Renders shamir step after PGP key chosen'
    );
  });

  test('it shows error with pgp key', async function (assert) {
    assert.expect(3);
    this.server.get('/sys/replication/dr/secondary/generate-operation-token/attempt', function () {
      return {};
    });
    this.server.post('/sys/replication/dr/secondary/generate-operation-token/attempt', () =>
      overrideResponse(400, { errors: ['error parsing PGP key'] })
    );
    await render(hbs`<Shamir::DrTokenFlow @action="generate-dr-operation-token" />`);
    await click(GENERAL.button('use-pgp-key-cta'));

    assert.ok(await waitUntil(() => find(SHAMIR_FORM.flowStep('choose-pgp-key'))), 'PGP form shows');
    await click(GENERAL.textToggle);
    await fillIn(GENERAL.textareaByAttr('pgp-key'), 'some-key-here');
    await click(GENERAL.button('use-pgp-key'));
    await click(GENERAL.submitButton);

    assert.ok(
      await waitUntil(() => find(SHAMIR_FORM.flowStep('primary-token'))),
      'shows primary token input after start'
    );
    await fillIn(GENERAL.inputByAttr('primary-token'), 'some-token');
    await click('[data-test-submit-primary-token]');

    await waitFor(GENERAL.messageError);
    assert.dom(GENERAL.messageError).hasText('Error error parsing PGP key');
  });

  test('it cancels correctly when generation not started', async function (assert) {
    assert.expect(2);
    const cancelSpy = sinon.spy();
    this.set('onCancel', cancelSpy);
    this.server.get('/sys/replication/dr/secondary/generate-operation-token/attempt', function () {
      return {};
    });
    this.server.delete('/sys/replication/dr/secondary/generate-operation-token/attempt', function () {
      assert.notOk('delete endpoint should not be queried');
      return {};
    });

    await render(
      hbs`<Shamir::DrTokenFlow @action="generate-dr-operation-token" @onCancel={{this.onCancel}} />`
    );
    assert.dom(GENERAL.cancelButton).hasText('Cancel', 'Close button has correct copy');
    await click(GENERAL.cancelButton);
    assert.ok(cancelSpy.calledOnce, 'cancel spy called on click');
  });

  test('it cancels correctly when generation has started but not finished', async function (assert) {
    assert.expect(6);
    const cancelSpy = sinon.spy(() => {
      assert.ok(true, 'passed cancel method called');
    });
    this.set('onCancel', cancelSpy);
    this.server.get('sys/replication/dr/secondary/generate-operation-token/attempt', function () {
      return {
        started: true,
        progress: 1,
        required: 3,
        complete: false,
      };
    });
    this.server.delete('/sys/replication/dr/secondary/generate-operation-token/attempt', function () {
      assert.ok(true, 'delete endpoint is queried');
      return {};
    });
    await render(
      hbs`<Shamir::DrTokenFlow @action="generate-dr-operation-token" @onCancel={{this.onCancel}} />`
    );

    await click(GENERAL.button('generate-token-cta'));
    assert.ok(
      await waitUntil(() => find(SHAMIR_FORM.flowStep('primary-token'))),
      'shows primary token step after start'
    );

    await fillIn(GENERAL.inputByAttr('primary-token'), 'some-token');
    await click('[data-test-submit-primary-token]');

    assert.dom(GENERAL.cancelButton).hasText('Cancel', 'Close button has correct copy');
    assert.ok(await waitUntil(() => find(GENERAL.inputByAttr('shamir-key'))), 'shows shamir key input');

    await click(GENERAL.cancelButton);
    assert.ok(
      await waitUntil(() => find(GENERAL.button('generate-token-cta'))),
      'shows generate token button again'
    );
  });

  test('it closes correctly when generation is completed', async function (assert) {
    assert.expect(2);
    const cancelSpy = sinon.spy();
    this.set('onCancel', cancelSpy);
    this.server.get('/sys/replication/dr/secondary/generate-operation-token/attempt', function () {
      return {
        started: true,
        progress: 3,
        required: 3,
        complete: true,
        encoded_token: 'foobar',
      };
    });
    this.server.delete('/sys/replication/dr/secondary/generate-operation-token/attempt', function () {
      assert.notOk('delete endpoint should not be queried');
      return {};
    });
    await render(
      hbs`<Shamir::DrTokenFlow @action="generate-dr-operation-token" @onCancel={{this.onCancel}} />`
    );

    await click(GENERAL.button('generate-token-cta'));
    await fillIn(GENERAL.inputByAttr('primary-token'), 'some-token');
    await click('[data-test-submit-primary-token]');

    await waitUntil(() => find(SHAMIR_FORM.flowStep('show-token')));
    assert.dom(GENERAL.cancelButton).hasText('Close', 'Close button has correct copy');
    await click(GENERAL.cancelButton);
    assert.ok(cancelSpy.calledOnce, 'cancel spy called on click');
  });
});
