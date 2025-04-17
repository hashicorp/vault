export interface AnalyticsProvider {
  name: string;
  identify: (identifier: string, traits: Record<string, string>) => void;
  start: (config: Record<string, string | boolean>) => void;
  trackPageView: (routeName: string, metadata: Record<string, string>) => void;
}
