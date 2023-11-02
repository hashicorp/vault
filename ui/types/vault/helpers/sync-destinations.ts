export type SyncDestinationType = 'aws-sm' | 'azure-kv' | 'gcp-sm' | 'gh' | 'vercel-project';
export type SyncDestinationName =
  | 'AWS Secrets Manager'
  | 'Azure Key Vault'
  | 'Google Secret Manager'
  | 'Github Actions'
  | 'Vercel Project';

export interface SyncDestination {
  name: SyncDestinationName;
  type: SyncDestinationType;
  icon: 'aws-color' | 'azure-color' | 'gcp-color' | 'github-color' | 'vercel-color';
  category: 'cloud' | 'dev-tools';
  maskedParams: Array<string>;
}
