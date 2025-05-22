/*
  Normally, the default export of this file would _do something_. In our case, the mthods for the provider simply log the 
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
