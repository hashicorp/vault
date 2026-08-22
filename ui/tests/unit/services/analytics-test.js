/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupTest } from 'vault/tests/helpers';
import { NAMESPACE } from 'vault/utils/preferences';
import { SegmentProvider } from 'vault/utils/analytics-providers/segment';

// Preference registry key backing telemetry consent (see app/utils/preferences.ts).
const CONSENT_KEY = `${NAMESPACE}:telemetryConsent`;
function setConsent(value) {
  window.localStorage.setItem(CONSENT_KEY, JSON.stringify(value));
}
function clearConsent() {
  window.localStorage.removeItem(CONSENT_KEY);
}

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

  // Analytics is non-essential. Verify a broken provider never surfaces to the user
  // or breaks navigation/interaction.
  module('provider error resilience', function (hooks) {
    hooks.afterEach(function () {
      sinon.restore();
    });

    test('start falls back to the no-op provider when the provider fails to initialize', function (assert) {
      sinon.stub(SegmentProvider.prototype, 'start').throws(new Error('sdk init failed'));

      this.service.start('segment', { enabled: true, write_key: 'test' });

      assert.strictEqual(this.service.provider.name, 'dummy', 'falls back to the no-op provider');
      assert.false(this.service.activated, 'analytics is not marked active');
    });

    test('trackEvent swallows provider errors and does not throw', function (assert) {
      this.service.provider = {
        name: 'boom',
        trackEvent: sinon.stub().throws(new Error('provider blew up')),
        trackPageView: sinon.stub(),
      };

      let error = null;
      try {
        this.service.trackEvent('button-click', { location: 'sidebar' });
      } catch (e) {
        error = e;
      }

      assert.strictEqual(error, null, 'the error is swallowed');
    });

    test('trackPageView swallows provider errors and does not throw', function (assert) {
      this.service.provider = {
        name: 'boom',
        trackEvent: sinon.stub(),
        trackPageView: sinon.stub().throws(new Error('provider blew up')),
      };

      let error = null;
      try {
        this.service.trackPageView('some.route');
      } catch (e) {
        error = e;
      }

      assert.strictEqual(error, null, 'the error is swallowed');
    });
  });

  module('Vault SM consent gate', function (hooks) {
    hooks.beforeEach(function () {
      clearConsent();

      sinon.stub(SegmentProvider.prototype, 'start');
      this.startSpy = sinon.spy(this.service, 'start');
      this.config = { enabled: true, write_key: 'test' };
    });

    hooks.afterEach(function () {
      sinon.restore();
      clearConsent();
    });

    module('#startVaultSmAnalytics', function () {
      test('operator disabled -> does not start and does not prompt', function (assert) {
        this.service.startVaultSmAnalytics(false, this.config);

        assert.true(this.startSpy.notCalled, 'analytics is not started');
        assert.false(this.service.shouldPromptConsent, 'no consent prompt');
        assert.false(this.service.operatorTelemetryEnabled, 'operator flag recorded as disabled');
      });

      test('operator enabled, consent unrecorded -> prompts and does not start', function (assert) {
        clearConsent();

        this.service.startVaultSmAnalytics(true, this.config);

        assert.true(this.service.shouldPromptConsent, 'consent banner is prompted');
        assert.true(this.startSpy.notCalled, 'analytics is not started yet');
      });

      test('operator enabled, consent = true -> starts with the segment provider', function (assert) {
        setConsent(true);

        this.service.startVaultSmAnalytics(true, this.config);

        assert.true(this.startSpy.calledOnceWith('segment', this.config), 'starts the segment provider');
        assert.false(this.service.shouldPromptConsent, 'no consent prompt');
      });

      test('operator enabled, consent = false -> does not start and does not prompt', function (assert) {
        setConsent(false);

        this.service.startVaultSmAnalytics(true, this.config);

        assert.true(this.startSpy.notCalled, 'analytics is not started');
        assert.false(this.service.shouldPromptConsent, 'no consent prompt');
      });

      test('localStorage unreadable -> does not start (treated as undecided)', function (assert) {
        sinon.stub(window.localStorage, 'getItem').throws(new Error('localStorage unavailable'));

        this.service.startVaultSmAnalytics(true, this.config);

        assert.true(this.startSpy.notCalled, 'analytics is not started when consent cannot be confirmed');
        assert.true(
          this.service.shouldPromptConsent,
          'prompts for consent as there was no readable decision'
        );
      });

      test('caches config and operator flag on every run for a later accept', function (assert) {
        setConsent(false);

        this.service.startVaultSmAnalytics(true, this.config);

        assert.strictEqual(this.service.pendingConfig, this.config, 'config is cached');
        assert.true(this.service.operatorTelemetryEnabled, 'operator flag is cached');
      });
    });

    module('#recordConsent', function () {
      test('accept when operator-enabled and not running -> persists true and starts', function (assert) {
        clearConsent();
        this.service.startVaultSmAnalytics(true, this.config); // prompts consent banner and caches config

        this.service.recordConsent(true);

        assert.strictEqual(window.localStorage.getItem(CONSENT_KEY), 'true', 'consent persisted as true');
        assert.true(this.startSpy.calledOnceWith('segment', this.config), 'starts in the current session');
        assert.false(this.service.shouldPromptConsent, 'banner dismissed');
      });

      test('accept when operator disabled -> persists true but does not start', function (assert) {
        this.service.startVaultSmAnalytics(false, this.config); // entry point is user preferences

        this.service.recordConsent(true);

        assert.strictEqual(window.localStorage.getItem(CONSENT_KEY), 'true', 'consent persisted');
        assert.true(this.startSpy.notCalled, 'does not start when operator has telemetry off');
      });

      test('decline while running -> persists false, resets and stops event emission', function (assert) {
        setConsent(true);
        this.service.startVaultSmAnalytics(true, this.config);
        const resetSpy = sinon.spy(this.service, 'reset');

        this.service.recordConsent(false);

        assert.strictEqual(window.localStorage.getItem(CONSENT_KEY), 'false', 'consent persisted as false');
        assert.true(resetSpy.calledOnce, 'reset is called to stop emission');
        assert.false(this.service.activated, 'analytics no longer active');
      });

      test('consent banner is not prompted if user preference was previously recorded', function (assert) {
        this.service.recordConsent(false);
        assert.notStrictEqual(window.localStorage.getItem(CONSENT_KEY), null, 'a decision is stored');

        // A stored decline suppresses the prompt on the next gate run.
        this.service.startVaultSmAnalytics(true, this.config);
        assert.false(this.service.shouldPromptConsent, 'recorded decline is not treated as undecided');
      });
    });

    module('#reset', function () {
      test('swaps provider out, removes the route listener, and clears activated', function (assert) {
        setConsent(true);
        this.service.startVaultSmAnalytics(true, this.config);
        const startedProvider = this.service.provider;
        const offSpy = sinon.spy(this.service.router, 'off');

        this.service.reset();

        assert.false(this.service.activated, 'no longer active');
        assert.notStrictEqual(
          this.service.provider,
          startedProvider,
          'provider swapped to the no-op provider'
        );
        assert.true(offSpy.called, 'route (page view) listener removed');
      });

      test('no-op when not activated', function (assert) {
        const provider = new ProviderStub();
        this.service.provider = provider;
        this.service.activated = false;

        this.service.reset();

        assert.strictEqual(this.service.provider, provider, 'provider unchanged when nothing is running');
      });
    });

    module('session transitions (no page reload)', function () {
      test('undecided -> accept -> analytics starts', function (assert) {
        clearConsent();
        this.service.startVaultSmAnalytics(true, this.config);

        this.service.recordConsent(true);

        assert.true(this.startSpy.calledOnceWith('segment', this.config), 'starts on accept');
      });

      test('undecided -> decline -> analytics does not start', function (assert) {
        clearConsent();
        this.service.startVaultSmAnalytics(true, this.config);

        this.service.recordConsent(false);

        assert.true(this.startSpy.notCalled, 'does not start on decline');
        assert.false(this.service.activated, 'not active');
      });

      test('decline -> accept via preferences -> starts in the same session', function (assert) {
        clearConsent();
        this.service.startVaultSmAnalytics(true, this.config);
        this.service.recordConsent(false);

        this.service.recordConsent(true); // later flips the preferences toggle on

        assert.true(
          this.startSpy.calledOnceWith('segment', this.config),
          'starts without a reload using the cached config'
        );
      });

      test('accept -> decline via preferences -> stops in the same session', function (assert) {
        clearConsent();
        this.service.startVaultSmAnalytics(true, this.config); // operator on, caches config, prompts
        this.service.recordConsent(true);
        const resetSpy = sinon.spy(this.service, 'reset');

        this.service.recordConsent(false);

        assert.true(resetSpy.calledOnce, 'stops emission on decline');
        assert.false(this.service.activated, 'no longer active');
      });

      test('accept -> reload -> decline -> stops', function (assert) {
        setConsent(true);
        this.service.startVaultSmAnalytics(true, this.config);

        this.service.recordConsent(false);

        assert.false(this.service.activated, 'declining after a reload stops analytics');
      });
    });
  });
});
