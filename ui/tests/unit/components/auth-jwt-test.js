/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import EmberObject from '@ember/object';
import Evented from '@ember/object/evented';
import sinon from 'sinon';
import { _cancelTimers as cancelTimers } from '@ember/runloop';

const mockWindow = EmberObject.extend(Evented, {
  origin: 'http://localhost:4200',
  close: () => {},
});

module('Unit | Component | auth-jwt', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.component = this.owner.lookup('component:auth-jwt');
    this.component.set('window', mockWindow.create());
    this.errorSpy = sinon.spy(this.component, 'handleOIDCError');
  });

  test('it should ignore messages from cross origin windows while waiting for oidc callback', async function (assert) {
    assert.expect(2);
    this.component.prepareForOIDC.perform(mockWindow.create());
    this.component.window.trigger('message', { origin: 'http://anotherdomain.com', isTrusted: true });

    assert.ok(this.errorSpy.notCalled, 'Error handler not triggered while waiting for oidc callback message');
    assert.strictEqual(this.component.exchangeOIDC.performCount, 0, 'exchangeOIDC method not fired');
    cancelTimers();
  });

  test('it should ignore untrusted messages while waiting for oidc callback', async function (assert) {
    assert.expect(2);
    this.component.prepareForOIDC.perform(mockWindow.create());
    this.component.window.trigger('message', { origin: 'http://localhost:4200', isTrusted: false });
    assert.ok(this.errorSpy.notCalled, 'Error handler not triggered while waiting for oidc callback message');
    assert.strictEqual(this.component.exchangeOIDC.performCount, 0, 'exchangeOIDC method not fired');
    cancelTimers();
  });

  // TODO: Flaky
  // test case for https://github.com/hashicorp/vault/issues/12436
  test('it should ignore messages sent from outside the app while waiting for oidc callback', async function (assert) {
    assert.expect(2);
    this.component.prepareForOIDC.perform(mockWindow.create());
    const message = {
      origin: 'http://localhost:4200',
      isTrusted: true,
      data: {
        namespace: 'foobar',
        path: '/foo/bar',
        state: 'authorized',
        code: 204,
      },
    };

    this.component.window.trigger('message', message);
    message.data.source = 'foo-bar';
    this.component.window.trigger('message', message);
    message.data.source = 'oidc-callback';
    this.component.window.trigger('message', message);

    assert.ok(this.errorSpy.notCalled, 'Error handler not triggered while waiting for oidc callback message');
    assert.strictEqual(
      this.component.exchangeOIDC.performCount,
      1,
      'exchangeOIDC method fires when oidc callback message is received'
    );
    cancelTimers();
  });
});
