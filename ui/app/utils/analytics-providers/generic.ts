export interface Provider {
  name: string;
  client: unknown;

  anonymousId: string;

  start(config: Record<string, string>): void;

  trackPageView(path: string, routeName: string, referrer: string): void;
  identifyUser(identifier: string, metadata: Record<string, string>): void;
}
