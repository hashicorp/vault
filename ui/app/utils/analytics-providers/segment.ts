/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { AnalyticsBrowser } from '@segment/analytics-next';

import type { MiddlewareFunction } from '@segment/analytics-next';
import type { AnalyticsEventName, AnalyticsProvider } from 'vault/vault/analytics';

interface SegmentConfig {
  enabled: boolean;
  write_key: string;
  isHvdManaged?: boolean;
}

// Allowlist of properties we intentionally send. Mirrors the PostHog redactEvent approach.
// Anything not listed here is stripped before the event leaves the browser.
//
// Values must be scalars. IBM explodes nested objects/arrays into distinct
// properties (e.g. `payload.policyType`) and has asked us not to rely on that,
// so event detail is carried in flat named properties such as `object`,
// `objectType`, `location` and `quantity` instead of a `payload` blob.
const ALLOWED_PROPERTIES = new Set([
  // IBM required instrumentation properties
  'productTitle',
  'productCode',
  'productCodeType',
  'UT30',
  'productPlanName',
  'productPlanType',
  'instanceId', // clusterId
  'subscriptionId', // licenseId, empty if community or HVD cluster
  'elementId',
  'namespace',
  'channel',
  'platformTitle',
  'action',
  'CTA',
  'location',
  'objectType',
  'object',
  'resultValue',
  'process',
  'successFlag',
  'quantity',
  'text',
  // CTA Clicked context
  'variation',
  'uiElement',
  'type',
  // Vault-specific
  'name',
  'routeName',
  'policy',
  // Browser context. Note: Segment auto-collects userAgent, timezone, screen and
  // library on the event `context` (not `properties`), so they reach Segment
  // regardless of this allowlist and do not need to be listed here.
  'locale',
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
      if (ALLOWED_PROPERTIES.has(key) || key.startsWith('custom.')) {
        redacted[key] = properties[key];
      }
    }
    payload.obj.properties = redacted;
  }

  // Segment's ingestion API backfills `context.ip` from the request's source IP
  // whenever the field is absent, so deleting it is not enough as the real IP
  // still lands in Segment. Setting it to a non-identifying placeholder tells
  // Segment there is nothing to infer. Also drop the auto-collected page context.
  payload.obj.context = payload.obj.context || {};
  payload.obj.context.ip = '0.0.0.0';
  delete payload.obj.context.page;

  next(payload);
};

// Static IBM instrumentation properties required on every event.
// instanceId, subscriptionId, and productPlanType are set dynamically at identify() time.
const IBM_STATIC_PROPERTIES = {
  productTitle: 'HASHICORP VAULT',
  productCode: '5621IJC',
  productCodeType: 'PID',
  UT30: '30GKT',
  platformTitle: 'HashiCorp Vault',
};

// productPlanType classifies the deployment model. Set dynamically at identify()
// from the isHvdManaged trait.
const PRODUCT_PLAN_TYPE = {
  HVD: 'Vault dedicated',
  SELF_MANAGED: 'Vault self-managed',
};

export const PROVIDER_NAME = 'segment';
export class SegmentProvider implements AnalyticsProvider {
  name = PROVIDER_NAME;

  client = new AnalyticsBrowser();
  licenseId = '';
  clusterId = '';
  userId = '';
  // Distinguishes enterprise from community clusters. Defaults to 'community'
  // until identify() runs; a cluster is only 'enterprise' once confirmed.
  productPlanName = 'community';
  // Deployment model (Vault dedicated vs self-managed).
  productPlanType = PRODUCT_PLAN_TYPE.SELF_MANAGED;
  instanceId = '';

  start(config: unknown) {
    const { enabled, write_key, isHvdManaged } = config as SegmentConfig;

    // Seed the deployment classification before any events fire.
    this.productPlanType = isHvdManaged ? PRODUCT_PLAN_TYPE.HVD : PRODUCT_PLAN_TYPE.SELF_MANAGED;

    if (enabled && write_key) {
      this.client.load({ writeKey: write_key });
      this.client.addSourceMiddleware(redactMiddleware);
    }
  }

  private get ibmProperties() {
    // Note: userId flows automatically via identify().
    // instanceId is the cluster ID - the deployed instance, unique per cluster
    // and present on both CE and Enterprise.
    // subscriptionId is the license ID - the commercial subscription. Community
    // clusters have no license, so it is omitted entirely (rather than sent
    // blank) to keep events schema-conformant.
    return {
      ...IBM_STATIC_PROPERTIES,
      productPlanName: this.productPlanName,
      productPlanType: this.productPlanType,
      ...(this.instanceId ? { instanceId: this.instanceId } : {}),
      ...(this.licenseId ? { subscriptionId: this.licenseId } : {}),
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

  identify(identifier: string, traits: Record<string, unknown>) {
    const isHvd = Boolean(traits['isHvdManaged']);
    this.userId = identifier;
    // HVD clusters have no customer subscription id available to the UI yet,
    // so omit subscriptionId for HVD rather than send the internal Vault license
    // id. Self-managed continues to use the license id.
    this.licenseId = isHvd ? '' : (traits['licenseId'] as string) || '';
    this.clusterId = (traits['clusterId'] as string) || '';
    this.instanceId = this.clusterId;
    this.productPlanName = traits['isEnterprise'] ? 'enterprise' : 'community';
    this.productPlanType = isHvd ? PRODUCT_PLAN_TYPE.HVD : PRODUCT_PLAN_TYPE.SELF_MANAGED;

    // Identify traits bypass the ALLOWED_PROPERTIES redaction (that only filters
    // track/page payloads), so for HVD strip the internal Vault license id from the
    // outbound traits.
    const outboundTraits = { ...traits };
    if (isHvd) {
      delete outboundTraits['licenseId'];
    }
    this.client.identify(identifier, { ...this.ibmProperties, ...outboundTraits });
  }

  trackPageView(routeName: string) {
    this.client.page(undefined, routeName, this.ibmProperties);
  }

  trackEvent(eventName: AnalyticsEventName, metadata?: Record<string, unknown>) {
    this.client.track(eventName, { ...this.ibmProperties, ...metadata });
  }
}
