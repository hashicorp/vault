/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

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
  PostHog collects a lot of stuff by default, this removes everything pretaining to urls and possible customer info in favor of generic information
*/
const redactEvent = (cr: CaptureResult | null) => {
  if (cr === null) return cr;
  const { properties } = cr;

  // extract ONLY what we need
  // this is definitely ridiculous, but it's the best way to be sure we're not accidentally including things that aren't wanted
  const {
    $browser,
    $browser_language,
    $browser_language_prefix,
    $browser_version,
    $device_id,
    $device_type,
    $groups,
    $insert_id,
    $is_identified,
    $lib,
    $lib_rate_limit_remaining_tokens,
    $lib_version,
    $os,
    $os_version,
    $pageview_id,
    $process_person_profile,
    $raw_user_agent,
    $recording_status,
    $screen_height,
    $screen_width,
    $sdk_debug_current_session_duration,
    $sdk_debug_replay_internal_buffer_length,
    $sdk_debug_replay_internal_buffer_size,
    $sdk_debug_retry_queue_size,
    $sdk_debug_session_start,
    $session_id,
    $time,
    $user_id,
    $viewport_height,
    $viewport_width,
    $window_id,
    distinct_id,
    routeName,
    token,
  } = properties;

  // replay into the sent object
  return {
    ...cr,
    properties: {
      $browser,
      $browser_version,
      $browser_language,
      $browser_language_prefix,
      $os,
      $os_version,
      $device_type,
      $raw_user_agent,
      $screen_height,
      $screen_width,
      $viewport_height,
      $viewport_width,
      $lib,
      $lib_version,
      $insert_id,
      $time,
      $device_id,
      $user_id,
      $groups,
      $session_id,
      $window_id,
      $recording_status,
      $sdk_debug_replay_internal_buffer_length,
      $sdk_debug_replay_internal_buffer_size,
      $sdk_debug_current_session_duration,
      $sdk_debug_session_start,
      $sdk_debug_retry_queue_size,
      $pageview_id,
      $is_identified,
      $process_person_profile,
      $lib_rate_limit_remaining_tokens,
      distinct_id,
      routeName,
      token,
    },
  };
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
        before_send: redactEvent,
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

  trackEvent(eventName: string, metadata: Record<string, string> = {}) {
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
