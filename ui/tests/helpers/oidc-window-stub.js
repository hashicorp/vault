/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
/* eslint-disable ember/no-computed-properties-in-native-classes */
import EmberObject, { computed } from '@ember/object';
import Evented from '@ember/object/evented';

export class NewFakeWindow {
  origin = 'https://my-vault.com';
  closed = false;

  #listeners = [];

  open() {
    this.closed = false;
  }

  close() {
    this.closed = true;
    this.#triggerClose();
  }

  on(event, callback) {
    if (event === 'close') {
      this.#listeners.push(callback);
    }
  }

  #triggerClose() {
    this.#listeners.forEach((callback) => callback());
  }
}

export const fakeWindow = EmberObject.extend(Evented, {
  init() {
    this._super(...arguments);
    this.on('close', () => {
      this.set('closed', true);
    });
  },
  screen: computed(function () {
    return {
      height: 600,
      width: 500,
    };
  }),
  origin: 'https://my-vault.com',
  closed: false,
  open() {},
  close() {},
});

export const buildMessage = (opts) => ({
  isTrusted: true,
  origin: 'https://my-vault.com',
  data: {
    source: 'oidc-callback',
    path: 'foo',
    state: 'state',
    code: 'code',
  },
  ...opts,
});
