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
import type { MountConfig, MountOptions } from 'vault/mount';

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
