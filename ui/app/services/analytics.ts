import Service, { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

import PostHogProvider, {
  PROVIDER_NAME as PostHogProviderName,
} from 'vault/utils/analytics-providers/posthog';
import GenericProvider, { PROVIDER_NAME as GenericProviderName } from 'vault/utils/analytics-providers/dummy';

import type Owner from '@ember/owner';
import type RouterService from '@ember/routing/router-service';
import type { Provider } from 'vault/utils/analytics-providers/generic';

export default class AnalyticsService extends Service {
  @service declare readonly router: RouterService;

  @tracked referrer = '';
  @tracked provider: Provider = new GenericProvider();

  constructor(owner: Owner) {
    super(owner);

    this.referrer = window.location.href;
    this.setupRouteEventListener();
  }

  private setupRouteEventListener() {
    // on successful route changes...
    this.router.on('routeDidChange', () => {
      const { currentRouteName, currentURL } = this.router;

      //
      this.trackPageView(currentURL || '', currentRouteName || '', this.referrer);
    });
  }

  trackPageView(currentRouteName: string, currentUrl: string, referrer: string) {
    this.provider.trackPageView(currentRouteName, currentUrl, referrer);
  }
  // tracking custom events will come later!

  // details
  start(provider: string, config: unknown) {
    if (!provider) return;

    this.provider = new PostHogProvider();

    switch (provider) {
      case GenericProviderName:
        this.provider = new GenericProvider();
        break;
      case PostHogProviderName:
        this.provider = new PostHogProvider();
        break;
    }

    this.provider.start(config);
  }

  identifyUser(identifier: string, metadata: Record<string, string>) {
    if (identifier) {
      console.log('I dunno, associate with that user?');
    }

    this.provider.identifyUser(identifier, metadata);
  }
}
