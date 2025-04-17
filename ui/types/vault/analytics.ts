export type AnalyticsProviderConfig = Record<string, string>;

export interface AnalyticsProvider {
  name: string;
  identify: (identifier: string, traits: Record<string, string>) => void;
  start: (config: AnalyticsProviderConfig) => void;
  trackPageView: (routeName: string, metadata: Record<string, string>) => void;
}
