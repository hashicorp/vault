/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';
import { encodeString, decodeString } from 'core/utils/b64';

module('Unit | Service | custom-messages', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    const payload = {
      keys: ['1', '2'],
      keyInfo: {
        1: {
          type: 'banner',
          message: encodeString('This is a banner message'),
          title: 'Look at banner Michael!',
          active: true,
        },
        2: {
          type: 'modal',
          message: encodeString('This is a modal message'),
          title: 'Otherwise known as a pop-up',
          active: true,
        },
      },
    };

    const api = this.owner.lookup('service:api');
    this.authMessagesApiStub = sinon
      .stub(api.sys, 'internalUiReadAuthenticatedActiveCustomMessages')
      .resolves(payload);
    this.unauthMessagesApiStub = sinon
      .stub(api.sys, 'internalUiReadUnauthenticatedActiveCustomMessages')
      .resolves(payload);

    this.customMessages = this.owner.lookup('service:custom-messages');

    this.messages = payload.keys.map((id) => {
      const data = payload.keyInfo[id];
      return {
        id,
        ...data,
        message: decodeString(data.message),
      };
    });

    this.authServiceStub = sinon.stub(this.owner.lookup('service:auth'), 'currentToken').value('token');
  });

  test('it should fetch unauthenticated messages', async function (assert) {
    this.authServiceStub.reset();
    await this.customMessages.fetchMessages();

    assert.true(this.unauthMessagesApiStub.called, 'API call made for unauthenticated messages');
  });

  test('it should fetch authenticated messages', async function (assert) {
    await this.customMessages.fetchMessages();

    assert.true(this.authMessagesApiStub.called, 'API call made for authenticated messages');
  });

  test('it should set messages from fetch response and decode message values', async function (assert) {
    await this.customMessages.fetchMessages();

    assert.deepEqual(
      this.customMessages.messages,
      this.messages,
      'messages are set correctly from API response'
    );
    assert.true(this.customMessages.bannerState['1'], 'banner state is set correctly');
  });

  test('it should filter banner messages', async function (assert) {
    await this.customMessages.fetchMessages();

    assert.deepEqual(
      this.customMessages.bannerMessages,
      [this.messages[0]],
      'only banner messages are returned'
    );
  });

  test('it should filter modal messages', async function (assert) {
    await this.customMessages.fetchMessages();

    assert.deepEqual(
      this.customMessages.modalMessages,
      [this.messages[1]],
      'only modal messages are returned'
    );
  });

  test('it should clear messages', async function (assert) {
    this.customMessages.messages = this.messages;

    this.customMessages.clearCustomMessages();
    assert.strictEqual(this.customMessages.messages.length, 0, 'messages are cleared');
  });

  test('it should clear messages on fetch error', async function (assert) {
    this.customMessages.messages = this.messages;
    this.authMessagesApiStub.rejects();

    await this.customMessages.fetchMessages();

    assert.strictEqual(this.customMessages.messages.length, 0, 'messages are cleared on fetch error');
  });

  test('it should update banner state on dismiss', async function (assert) {
    this.customMessages.onBannerDismiss('1');
    assert.false(this.customMessages.bannerState['1'], 'banner state is updated correctly on dismiss');
  });
});
