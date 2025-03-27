import type { ProviderConfig, Provider } from 'vault/utils/analytics-providers/generic';


class DummyClient {
  identify() {
    // noop
  }

  capture(eventName: string, metadata: Record<string, string> = {}, options = {}) {
    console.log(eventName, metadata, options);
  }
}

export const PROVIDER_NAME = 'dummy';

export default class DummyProvider implements Provider {
  name = PROVIDER_NAME;

  client: unknown = null;

  get anonymousId() {
    return crypto.randomUUID();
  }

  start(_c: ProviderConfig) {
    console.log(`[Analytics] :: Starting ${this.name} provider`);
    this.client = new DummyClient();
  }

  trackPageView(path: string, routeName: string, referrer: string): void {
    console.log(`[Analytics] :: trackPageView :: ${routeName}`);
    this.client.capture('$pageview', { routeName });
  }

  identifyUser(identifier, metadata) {
    if (identifier) {
      console.log(`[Analytics] :: identifyUser (${identifier})::`, metadata);
    } else {
      console.log(`[Analytics] :: identifyUser (anonymous)::`, metadata);
    }
  }
}
