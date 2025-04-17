import posthog from 'posthog-js/dist/module.no-external';

import type { AnalyticsProvider } from 'vault/vault/analytics';
import type { CaptureResult } from 'posthog-js/dist/module.no-external';

interface PostHogConfig {
  enabled: boolean;
  project_id: string;
  api_host: string;
}

export const PROVIDER_NAME = 'dummy';

export class PostHogProvider implements AnalyticsProvider {
  name = PROVIDER_NAME;
  client = posthog;

  get anonymousId() {
    return this.client.get_distinct_id();
  }

  start(config: unknown) {
    const { enabled, project_id, api_host } = config as PostHogConfig;

    if (enabled) {
      this.client.init(project_id, {
        api_host,
        person_profiles: 'identified_only',
        persistence: 'memory',
        autocapture: false,
        disable_session_recording: true,
        advanced_disable_decide: true,
        capture_pageview: false,
        // replace with real redact operation
        before_send: (cr: CaptureResult | null) => cr,
      });
    }
  }

  identify(identifier: string, traits: Record<string, string>) {
    this.client.identify(identifier, {
      ...traits,
    });
  }

  trackPageView(routeName: string /* metadata: Record<string, string> */) {
    this.client.capture('$pageview', {
      routeName,
    });
  }
}
