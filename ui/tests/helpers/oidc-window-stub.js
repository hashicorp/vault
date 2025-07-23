/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import sinon from 'sinon';

// suggestions for a custom popup
// passing { close: true } automatically closes popups opened from window.open()
// passing { closed: true } sets value on popup window
export const windowStub = ({ stub, popup } = {}) => {
  // if already stubbed, don't re-stub
  const openStub = stub ? stub : sinon.stub(window, 'open');

  const defaultPopup = { close: () => true };
  openStub.returns(popup || defaultPopup);
  return openStub;
};

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
