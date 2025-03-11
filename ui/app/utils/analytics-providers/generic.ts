export type ProviderConfig = {
  API_KEY: string;
  api_host: string;
  enabled: boolean;
}

export interface Provider {
  name: string;
  client: unknown;

  anonymousId: string;

  start(config: ProviderConfig): void;

  trackPageView(path: string, routeName: string, referrer: string): void;
  identifyUser(identifier: string, metadata: Record<string, string>): void;
}
