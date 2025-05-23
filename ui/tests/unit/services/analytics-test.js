/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';

import sinon from 'sinon';

import { setupTest } from 'vault/tests/helpers';

class ProviderStub {
  name = 'testing';

  start = sinon.stub();
  identify = sinon.stub();
  trackPageView = sinon.stub();
}

module('Unit | Service | analytics', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.service = this.owner.lookup('service:analytics');
  });

  hooks.afterEach(function () {
    sinon.reset();
  });

  test('#identifyUser passes data to the provider', function (assert) {
    const providerStub = new ProviderStub();
    this.service.provider = providerStub;

    const identifier = 'carl';
    const traits = { apples: 'oranges' };

    this.service.identifyUser(identifier, traits);

    assert.true(providerStub.identify.calledOnce, 'the service calls identify on the provider');
    assert.true(
      providerStub.identify.calledWith(identifier, traits),
      'the provider recieves the expected id and traits'
    );
  });

  test('#trackPageView passes data to the provider', function (assert) {
    const providerStub = new ProviderStub();
    this.service.provider = providerStub;

    this.service.trackPageView('test', { currentRouteName: 'ham' });

    assert.true(providerStub.trackPageView.called, 'it calls the tracking method on the provider');
    assert.true(
      providerStub.trackPageView.calledWith('test', { currentRouteName: 'ham' }),
      'it passes the correct args to the provider'
    );
  });

  module('#log', function (hooks) {
    hooks.beforeEach(function () {
      // eslint-disable-next-line no-console
      console.log = sinon.stub(console, 'log');
    });

    hooks.afterEach(function () {
      // eslint-disable-next-line no-console
      console.log.restore();
    });

    test('logging is not shown when inactive', function (assert) {
      this.service.debug = false;
      // for the next few lines, console.log WILL NOT WORK AS EXPECTED
      this.service.trackPageView('a', null);

      // eslint-disable-next-line no-console
      assert.true(console.log.notCalled, 'console.log is called');
    });
  });
});
