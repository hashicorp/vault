import Service, { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

import { DummyProvider, PROVIDER_NAME as DummyProviderName } from 'vault/utils/analytics-providers/dummy';

import type { AnalyticsProvider } from 'vault/vault/analytics';
import type Owner from '@ember/owner';
import type RouterService from '@ember/routing/router-service';

export default class AnalyticsService extends Service {
  @service declare readonly router: RouterService;

  @tracked provider: AnalyticsProvider = new DummyProvider();

  debug = true;

  private log(...args: unknown[]) {
    if (this.debug) {
      // eslint-disable-next-line no-console
      console.log(`[Analytics - ${this.provider.name}]`, ...args);
    }
  }

  private setupRouteEventListener() {
    // on successful route changes...
    this.router.on('routeDidChange', () => {
      const { currentRouteName, currentURL } = this.router;

      this.trackPageView(currentURL || '', {
        currentRouteName: currentRouteName || '',
      });
    });
  }

  constructor(owner: Owner) {
    super(owner);

    this.setupRouteEventListener();
  }

  identifyUser(identifer: string, traits: Record<string, string>) {
    this.provider.identify(identifer, traits);
    this.log('identifyUser', identifer, traits);
  }

  start(provider: string, config = {}) {
    // fail silently, analytics is nonessential
    if (!provider) return;

    switch (provider) {
      case DummyProviderName:
        this.provider = new DummyProvider();
        break;
    }

    this.provider.start(config);

    this.log('start');
  }

  trackPageView(routeName: string, metadata: Record<string, string>) {
    this.provider.trackPageView(routeName, metadata);

    this.log('$pageview', routeName, metadata);
  }
}
