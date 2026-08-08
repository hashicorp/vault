/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { AnalyticsBrowser } from '@segment/analytics-next';

import type { AnalyticsProvider } from 'vault/vault/analytics';
import type { MiddlewareFunction } from '@segment/analytics-next';

interface SegmentConfig {
  enabled: boolean;
  write_key: string;
}

// Allowlist of properties we intentionally send. Mirrors the PostHog redactEvent approach.
// Anything not listed here is stripped before the event leaves the browser.
const ALLOWED_PROPERTIES = new Set([
  // IBM required instrumentation properties
  'productTitle',
  'productCode',
  'productCodeType',
  'UT30',
  'productPlanName',
  'productPlanType',
  'instanceId',
  'subscriptionId',
  'elementId',
  'namespace',
  'channel',
  'platformTitle',
  'action',
  'CTA',
  'location',
  'objectType',
  'process',
  'payload',
  'text',
  // Vault-specific
  'name',
  'routeName',
  'policy',
  // Browser context
  'userAgent',
  'locale',
  'timezone',
  'screen',
  'library',
  // Viewport / form-factor (mobile vs desktop)
  'viewportWidth',
  'viewportHeight',
  'viewportOrientation',
]);

const redactMiddleware: MiddlewareFunction = ({ payload, next }) => {
  const properties = payload.obj.properties as Record<string, unknown> | undefined;

  if (properties) {
    const redacted: Record<string, unknown> = {};
    for (const key of Object.keys(properties)) {
      // Keep allowlisted properties plus any `custom.*` property, which IBM's
      // Tracking Plan permits on every event type via its `^custom\.` pattern.
      // This is how Created Object carries its descriptive detail, since that
      // schema has no `payload` property of its own.
      if (ALLOWED_PROPERTIES.has(key) || key.startsWith('custom.')) {
        redacted[key] = properties[key];
      }
    }
    payload.obj.properties = redacted;
  }

  // Remove page-level context properties that are automatically added by Segment and strip out IP address
  if (payload.obj.context) {
    delete payload.obj.context.ip;
    delete payload.obj.context.page;
  }

  next(payload);
};

// Static IBM instrumentation properties required on every event.
// instanceId and subscriptionId are dynamically set from licenseId at identify() time.
const IBM_STATIC_PROPERTIES = {
  productTitle: 'HASHICORP VAULT',
  productCode: '5621IJC',
  productCodeType: 'PID',
  UT30: '30GKT',
  productPlanName: 'internal',
  productPlanType: 'internal',
  platformTitle: 'Vault',
};

// Map Vault-specific descriptive event names (sent as-is to PostHog for HVD) to
// the IBM Tracking Plan's generic event names for Segment. Event names not listed
// here are already IBM-conformant and pass through unchanged.
// Once HVD is fully instrumented, we can remove this mapping and send the IBM
// event names directly.
const IBM_EVENT_NAME_MAP: Record<string, string> = {
  'vault_ui_core_web-repl_toggle': 'UI Interaction',
};

export const PROVIDER_NAME = 'segment';

export class SegmentProvider implements AnalyticsProvider {
  name = PROVIDER_NAME;

  client = new AnalyticsBrowser();
  licenseId = '';
  userId = '';
  // instanceId and subscriptionId are both the licenseId per IBM Instrument Growth guidance
  instanceId = '';

  start(config: unknown) {
    const { enabled, write_key } = config as SegmentConfig;

    if (enabled && write_key) {
      this.client.load({ writeKey: write_key });
      this.client.addSourceMiddleware(redactMiddleware);
    }
  }

  private get ibmProperties() {
    // Note: userId flows automatically via identify(), and the license is carried
    // by instanceId/subscriptionId. Sending raw userId/licenseId as top-level event
    // properties is redundant and rejected by IBM's Tracking Plan schema.
    return {
      ...IBM_STATIC_PROPERTIES,
      ...(this.instanceId ? { instanceId: this.instanceId, subscriptionId: this.instanceId } : {}),
      ...this.viewportProperties,
    };
  }

  // Computed values for mobile-vs-desktop analysis. These are raw viewport
  // measurements (not device identifiers).
  // We omit IBM's `viewportWidthGroup` (its buckets follow Carbon breakpoints
  // whereas Vault uses HDS)
  private get viewportProperties() {
    if (typeof window === 'undefined') return {};

    const width = window.innerWidth;
    const height = window.innerHeight;

    return {
      viewportWidth: width,
      viewportHeight: height,
      viewportOrientation: width >= height ? 'landscape-primary' : 'portrait-primary',
    };
  }

  identify(identifier: string, traits: Record<string, string>) {
    this.userId = identifier;
    this.licenseId = traits['licenseId'] || '';
    this.instanceId = this.licenseId;
    this.client.identify(identifier, { ...this.ibmProperties, ...traits });
  }

  trackPageView(routeName: string) {
    this.client.page(undefined, routeName, this.ibmProperties);
  }

  trackEvent(eventName: string, metadata?: Record<string, unknown>) {
    const ibmEventName = IBM_EVENT_NAME_MAP[eventName] ?? eventName;
    this.client.track(ibmEventName, { ...this.ibmProperties, ...metadata });
  }
}
