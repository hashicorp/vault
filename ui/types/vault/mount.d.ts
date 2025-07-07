export type MountConfig = {
  forceNoCache?: boolean;
  listingVisibility?: string | boolean;
  defaultLeaseTtl?: number;
  maxLeaseTtl?: number;
  allowedManagedKeys?: string[];
  auditNonHmacRequestKeys?: string[];
  auditNonHmacResponseKeys?: string[];
  passthroughRequestHeaders?: string[];
  allowedResponseHeaders?: string[];
  identityTokenKey?: string;
};

export type MountOptions = {
  version: number;
};

export type Mount = {
  path: string;
  accessor: string;
  config: MountConfig;
  description: string;
  externalEntropyAccess: boolean;
  local: boolean;
  options?: MountOptions;
  pluginVersion: string;
  runningPluginVersion: string;
  runningSha256: string;
  sealWrap: boolean;
  type: string;
  uuid: string;
};
