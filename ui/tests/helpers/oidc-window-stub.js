/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import EmberObject from '@ember/object';
import Evented from '@ember/object/evented';

export class WindowStub extends EventTarget {
  close() {
    this.dispatchEvent(new CustomEvent('close')); // Trigger 'close' event using CustomEvent
  }
}

// using Evented is deprecated, but it's the only way we can trigger a message that is trusted
// by calling window.trigger. Using dispatchEvent will always result in an untrusted event.
export const fakeWindow = EmberObject.extend(Evented, {
  init() {
    this._super(...arguments);
    this.on('close', () => {
      this.set('closed', true);
    });
  },
  get screen() {
    return {
      height: 600,
      width: 500,
    };
  },
  origin: 'https://my-vault.com',
  closed: false,
  open() {},
  close() {},
});

export const buildMessage = (opts) => ({
  isTrusted: true,
  origin: 'https://my-vault.com',
  data: callbackData(),
  ...opts,
});

export const callbackData = (data = {}) => ({
  source: 'oidc-callback',
  path: 'foo',
  state: 'state',
  code: 'code',
  ...data,
});
