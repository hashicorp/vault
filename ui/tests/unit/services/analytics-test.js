/**
 * Copyright IBM Corp. 2016, 2025
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
      sinon.stub(console, 'info');
    });

    hooks.afterEach(function () {
      console.info.restore();
    });

    test('logging does not show outside of dev environment', function (assert) {
      this.service.debug = false;
      this.service.trackPageView('test-route', { foo: 'bar' });

      assert.true(console.info.notCalled, 'console.info is not called when debug is false');
    });

    test('logging shows in dev environments with correct format', function (assert) {
      this.service.debug = true;
      this.service.trackPageView('test-route', { foo: 'bar' });

      assert.true(
        console.info.calledOnceWith('[Analytics - dummy]', '$pageview', 'test-route', { foo: 'bar' }),
        'console.info is called once with correctly formatted message'
      );
    });

    test('logging works for all public methods', function (assert) {
      this.service.debug = true;

      this.service.identifyUser('user-123', { role: 'admin' });
      this.service.trackEvent('button-click', { location: 'sidebar' });

      assert.strictEqual(console.info.callCount, 2, 'log is called for each public method');
    });
  });
});
