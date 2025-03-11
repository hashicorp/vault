// No external code loading possible (this disables all extensions such as Replay, Surveys, Exceptions etc.)
import posthog from 'posthog-js/dist/module.no-external';

import type { ProviderConfig, Provider } from 'vault/utils/analytics-providers/generic';

export const PROVIDER_NAME = 'PostHog';

export default class PostHogProvider implements Provider {
  name = PROVIDER_NAME;
  client = posthog;

  get anonymousId() {
    return this.client.get_distinct_id();
  }

  start(config: ProviderConfig) {
    const { enabled, API_KEY, api_host } = config;

    if (enabled) {
      console.log(`[tracking] - start ${this.name} tracking service`);
      // this.client.init(API_KEY, {
      //   api_host,
      //   person_profiles: 'identified_only',
      //   persistence: 'memory',
      //   autocapture: false,
      //   disable_session_recording: true,
      //   advanced_disable_decide: true,
      //   capture_pageview: false,
      //   before_send: redact,
      // });
    }
  }

  trackPageView(path: string, currentRouteName: string, referrer: string) {
    this.client.capture('$pageview', {
      currentRouteName,
      referrer,
    });
  }

  identifyUser() {
    console.log('[tracking] - IDENTIFYING USER');
  }
}