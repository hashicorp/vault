import posthog from 'posthog-js/dist/module.no-external';

import type { AnalyticsProvider } from 'vault/vault/analytics';
import type { CaptureResult } from 'posthog-js/dist/module.no-external';

interface PostHogConfig {
  enabled: boolean;
  project_id: string;
  api_host: string;
}

/*
  formatEvent takes the default posthog capture event and trims it down to only the things we want to collect
  PostHog collects a lot of stuff by default, this removes most of those in favor of generic information
*/
const formatEvent = (cr: CaptureResult | null) => {
  if (cr === null) return cr;
  return cr;
};

export const PROVIDER_NAME = 'posthog';

export class PostHogProvider implements AnalyticsProvider {
  name = PROVIDER_NAME;

  client = posthog;
  licenseId = '';

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
        before_send: formatEvent,
      });
    }
  }

  identify(identifier: string, traits: Record<string, string>) {
    this.licenseId = traits['licenseId'] || '';

    this.client.identify(identifier, {
      ...traits,
    });
  }

  trackPageView(routeName: string /* metadata: Record<string, string> */) {
    // use licenseId as a grouping for this cluster
    if (this.licenseId) {
      this.client.capture('$pageview', {
        routeName,
        $groups: {
          licenseId: this.licenseId,
        },
      });
    } else {
      this.client.capture('$pageview', { routeName });
    }
  }

  trackEvent(eventName: string, metadata: Record<string, string>) {
    // use licenseId as a grouping for this cluster
    if (this.licenseId) {
      this.client.capture(eventName, {
        ...metadata,
        $groups: {
          licenseId: this.licenseId,
        },
      });
    } else {
      this.client.capture(eventName, metadata);
    }
  }
}
