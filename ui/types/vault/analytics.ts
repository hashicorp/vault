/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export interface AnalyticsProvider {
  name: string;
  identify: (identifier: string, traits: Record<string, string>) => void;
  start: (config: Record<string, string | boolean>) => void;
  trackPageView: (routeName: string, metadata: Record<string, string>) => void;
  trackEvent: (eventName: string, metadata: Record<string, string>) => void;
}

export interface AnalyticsConfig extends Record<string, string | boolean> {
  enabled: boolean;
}
