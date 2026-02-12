/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

import { DummyProvider, PROVIDER_NAME as DummyProviderName } from 'vault/utils/analytics-providers/dummy';
import {
  PostHogProvider,
  PROVIDER_NAME as PostHogProviderName,
} from 'vault/utils/analytics-providers/posthog';

import type { AnalyticsConfig, AnalyticsProvider } from 'vault/vault/analytics';
import type RouterService from '@ember/routing/router-service';

import config from 'vault/config/environment';

export default class AnalyticsService extends Service {
  @service declare readonly router: RouterService;

  @tracked activated = false;
  @tracked provider: AnalyticsProvider = new DummyProvider();

  debug = config.environment === 'development';

  private log(...args: unknown[]) {
    if (this.debug) {
      console.info(`[Analytics - ${this.provider.name}]`, ...args);
    }
  }

  private setupRouteEventListener() {
    // on successful route changes...
    this.router.on('routeDidChange', () => {
      const { currentRouteName } = this.router;

      this.trackPageView(currentRouteName || 'unknown-route');
    });
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
      }

      // only start things once we've confirmed we want to
      this.provider.start(config);
      this.activated = true;
      this.setupRouteEventListener();

      this.log('start');
    }
  };

  trackPageView = (routeName: string, metadata?: Record<string, string>) => {
    this.provider.trackPageView(routeName, metadata || {});

    this.log('$pageview', routeName, metadata);
  };

  trackEvent = (eventName: string, metadata: Record<string, string>) => {
    this.provider.trackEvent(eventName, metadata);

    this.log('custom event', eventName, metadata);
  };
}
