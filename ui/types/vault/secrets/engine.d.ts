/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type {
  MountsEnableSecretsEngineRequest,
  AwsConfigureRootIamCredentialsRequest,
  AwsConfigureLeaseRequest,
  AzureConfigureRequest,
  GoogleCloudConfigureRequest,
} from '@hashicorp/vault-client-typescript';

export type EngineConfig = {
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

export type EngineOptions = {
  version: number;
};

export type SecretsEngine = {
  path: string;
  accessor: string;
  config: EngineConfig;
  description: string;
  externalEntropyAccess: boolean;
  local: boolean;
  options?: EngineOptions;
  pluginVersion: string;
  runningPluginVersion: string;
  runningSha256: string;
  sealWrap: boolean;
  type: string;
  uuid: string;
};

type CommonConfigParams = {
  rotationPeriod: number;
  rotationSchedule: string;
  rotationWindow: number;
  identityTokenAudience: string;
  identityTokenTtl: number;
  disableAutomatedRotation: boolean;
  issuer?: string;
};

export type AwsConfig = CommonConfigParams & {
  accessKey: string;
  iamEndpoint: string;
  maxRetries: number;
  region: string;
  roleArn: string;
  stsEndpoint: string;
  stsFallbackEndpoints: string[];
  stsFallbackRegions: string[];
  stsRegion: string;
  usernameTemplate: string;
  lease?: string;
  leaseMax?: string;
};

export type AzureConfig = CommonConfigParams & {
  clientId: string;
  environment: string;
  subscriptionId: string;
  tenantId: string;
};

export type GcpConfig = CommonConfigParams & {
  maxTtl: number;
  serviceAccountEmail: string;
  ttl: number;
};

export type SshConfig = {
  publicKey: string;
  generateSigningKey: boolean;
};

export type SecretsEngineFormData = MountsEnableSecretsEngineRequest & {
  path: string;
  config?: EngineConfig;
  options?: EngineOptions;
  kvConfig?: {
    maxVersions?: number;
    casRequired?: boolean;
    deleteVersionAfter?: string;
  };
};

type Issuer = {
  issuer?: string;
};

export type AwsConfigFormData = AwsConfigureRootIamCredentialsRequest & AwsConfigureLeaseRequest & Issuer;
export type AzureConfigFormData = AzureConfigureRequest & Issuer;
export type GcpConfigFormData = GoogleCloudConfigureRequest & Issuer;
