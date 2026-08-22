/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

import { DummyProvider, PROVIDER_NAME as DummyProviderName } from 'vault/utils/analytics-providers/dummy';
import {
  PostHogProvider,
  PROVIDER_NAME as PostHogProviderName,
} from 'vault/utils/analytics-providers/posthog';

import {
  SegmentProvider,
  PROVIDER_NAME as SegmentProviderName,
} from 'vault/utils/analytics-providers/segment';

import { getPreference, hasPreference, setPreference } from 'vault/utils/preferences';

import type { AnalyticsConfig, AnalyticsProvider } from 'vault/vault/analytics';
import type RouterService from '@ember/routing/router-service';

import config from 'vault/config/environment';

export default class AnalyticsService extends Service {
  @service declare readonly router: RouterService;

  @tracked activated = false;
  @tracked provider: AnalyticsProvider = new DummyProvider();
  @tracked shouldPromptConsent = false;

  // Cached by the Vault SM gate so a later accept (banner or Preferences) can
  // start analytics in the current session without a page reload, and so the
  // operator opt-in is re-checked at accept time.
  private pendingConfig?: AnalyticsConfig;
  private operatorTelemetryEnabled = false;

  debug = config.environment === 'development';

  private log(...args: unknown[]) {
    if (this.debug) {
      console.info(`[Analytics - ${this.provider.name}]`, ...args);
    }
  }

  private routeDidChange = () => {
    const { currentRouteName } = this.router;
    this.trackPageView(currentRouteName || 'unknown-route');
  };

  private setupRouteEventListener() {
    this.router.on('routeDidChange', this.routeDidChange);
  }

  private teardownRouteEventListener() {
    this.router.off('routeDidChange', this.routeDidChange);
  }

  identifyUser = (identifer: string, traits: Record<string, string>) => {
    this.provider.identify(identifer, traits);
    this.log('identifyUser', identifer, traits);
  };

  start = (provider: string, config: AnalyticsConfig) => {
    // fail silently, analytics is nonessential
    if (!provider) {
      this.provider = new DummyProvider();
      this.debug = false;
      return;
    }

    // if analytics are not enabled, don't start the service
    if (config.enabled) {
      switch (provider) {
        case DummyProviderName:
          this.provider = new DummyProvider();
          break;
        case PostHogProviderName:
          this.provider = new PostHogProvider();
          break;
        case SegmentProviderName:
          this.provider = new SegmentProvider();
          break;
      }

      // Only start things once we've confirmed we want to. If the provider
      // fails to initialize then fall back to the no-op provider so analytics
      // stays silent and non-essential.
      try {
        this.provider.start(config);
        this.activated = true;
        this.setupRouteEventListener();
        this.log('start');
      } catch (e) {
        this.provider = new DummyProvider();
        this.activated = false;
        this.log('provider failed to start, falling back to no-op', e);
      }
    }
  };

  /**
   * Decision matrix (operator flag AND per-user consent must both allow it):
   *   operator off                     -> do nothing
   *   operator on, consent unrecorded  -> prompt (banner) but do not start
   *   operator on, consent = true      -> start analytics
   *   operator on, consent = false     -> do nothing, do not prompt (avoid noise)
   */
  startVaultSmAnalytics = (uiTelemetryEnabled: boolean, config: AnalyticsConfig) => {
    this.shouldPromptConsent = false;
    this.operatorTelemetryEnabled = uiTelemetryEnabled;

    if (!uiTelemetryEnabled) return;

    // Cache the config on every gate run (not just when prompting) so a later
    // accept via user preferences can start in the same session even after a reload.
    this.pendingConfig = config;

    if (this.activated) return;

    let consentRecorded: boolean;
    let consentGranted: boolean;
    try {
      consentRecorded = hasPreference('telemetryConsent');
      consentGranted = getPreference('telemetryConsent');
    } catch (e) {
      // localStorage unavailable (private browsing, quota) then fail safe.
      this.log('consent unreadable, analytics not started');
      return;
    }

    // User has never decided so surface the consent banner to record consent.
    if (!consentRecorded) {
      this.shouldPromptConsent = true;
      this.log('consent unrecorded, prompting');
      return;
    }

    if (consentGranted) {
      this.start(SegmentProviderName, config);
    }
  };

  /**
   * Records the user's telemetry consent decision. Both acceptance and decline
   * are persisted (not just acceptance) so that a decline is distinguishable from
   * never decided. This is to ensure the banner is not noisy
   *
   *   accept  -> start analytics now (if operator-enabled and not already running)
   *   decline -> stop analytics now (if running), ceasing event emission this session
   */
  recordConsent = (accepted: boolean) => {
    setPreference('telemetryConsent', accepted);
    this.shouldPromptConsent = false;
    this.log('consent recorded', accepted);

    if (accepted) {
      // Re-check the operator opt-in at accept time. The Preferences toggle can
      // fire outside the gate, and the operator flag must still allow telemetry.
      if (this.operatorTelemetryEnabled && this.pendingConfig && !this.activated) {
        this.start(SegmentProviderName, this.pendingConfig);
      }
    } else if (this.activated) {
      this.reset();
    }
  };

  /**
   * Stops analytics for the current session. This swaps the provider back to the
   * no-op DummyProvider, removes the page-view route listener, and clears the
   * activated flag. Used when the user revokes consent so that Vault SM
   * telemetry ceases without reload.
   */
  reset = () => {
    if (!this.activated) return;
    this.teardownRouteEventListener();
    this.provider = new DummyProvider();
    this.activated = false;
    this.log('reset');
  };

  // Swallow provider errors as analytics is non-essential and
  // must never break navigation or interaction if the provider throws.
  trackPageView = (routeName: string, metadata?: Record<string, string>) => {
    try {
      this.provider.trackPageView(routeName, metadata || {});
      this.log('$pageview', routeName, metadata);
    } catch (e) {
      this.log('trackPageView failed', e);
    }
  };

  // Swallow provider errors as analytics is non-essential and
  // must never break navigation or interaction if the provider throws.
  trackEvent = (eventName: string, metadata: Record<string, unknown>) => {
    try {
      this.provider.trackEvent(eventName, metadata);
      this.log('custom event', eventName, metadata);
    } catch (e) {
      this.log('trackEvent failed', e);
    }
  };
}
