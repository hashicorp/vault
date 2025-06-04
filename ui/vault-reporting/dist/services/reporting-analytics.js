import { getOwner } from '@ember/owner';
import Service from '@ember/service';

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * This service is used to look up the `analytics` service in the host application and track events if it exists. If it doesn't exist
 * or the implementation breaks it falls back gracefully to do nothing.
 */
class ReportingAnalytics extends Service {
  // Using the `@service` decorator will throw an error if the service  does not exist on the host application.
  // This allows us to be defensive and have `trackEvent` be a no-op if the service is not present.
  get analytics() {
    return getOwner(this)?.lookup('service:analytics');
  }
  trackEvent(event, properties, options) {
    if (!this.analytics?.trackEvent) {
      return;
    }
    try {
      const prefix = 'vault_reporting';
      const prefixedEvent = `${prefix}_${event}`;
      this.analytics.trackEvent(prefixedEvent, properties, options);
    } catch (e) {
      // no-op
      console.warn('Error tracking event:', e);
    }
  }
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.

export { ReportingAnalytics as default };
//# sourceMappingURL=reporting-analytics.js.map
