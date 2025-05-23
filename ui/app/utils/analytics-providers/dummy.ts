/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/*
  Normally, the default export of this file would _do something_. 
  In the dummy case, the methods for the provider are noops so we
    have a safe fallback when analytics is disabled.
*/
import type { AnalyticsProvider } from 'vault/vault/analytics';

export const PROVIDER_NAME = 'dummy';

export class DummyProvider implements AnalyticsProvider {
  name = PROVIDER_NAME;

  start() {
    /* intentionally blank */
  }

  identify() {
    /* intentionally blank */
  }

  trackPageView() {
    /* intentionally blank */
  }

  trackEvent() {
    /* intentionally blank */
  }
}
