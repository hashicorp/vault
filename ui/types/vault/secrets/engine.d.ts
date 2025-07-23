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
  rotation_period: number;
  rotation_schedule: string;
  rotation_window: number;
  identity_token_audience: string;
  identity_token_ttl: number;
  disable_automated_rotation: boolean;
  issuer?: string;
};

export type AwsConfig = CommonConfigParams & {
  access_key: string;
  iam_endpoint: string;
  max_retries: number;
  region: string;
  role_arn: string;
  sts_endpoint: string;
  sts_fallback_endpoints: string[];
  sts_fallback_regions: string[];
  sts_region: string;
  username_template: string;
  lease?: string;
  lease_max?: string;
};

export type AzureConfig = CommonConfigParams & {
  client_id: string;
  environment: string;
  subscription_id: string;
  tenant_id: string;
};

export type GcpConfig = CommonConfigParams & {
  max_ttl: number;
  service_account_email: string;
  ttl: number;
};

export type SshConfig = {
  public_key: string;
  generate_signing_key: boolean;
};

export type SecretsEngineFormData = MountsEnableSecretsEngineRequest & {
  path: string;
  config?: MountConfig;
  options?: MountOptions;
  kv_config?: {
    max_versions?: number;
    cas_required?: boolean;
    delete_version_after?: string;
  };
};

type Issuer = {
  issuer?: string;
};

export type AwsConfigFormData = AwsConfigureRootIamCredentialsRequest & AwsConfigureLeaseRequest & Issuer;
export type AzureConfigFormData = AzureConfigureRequest & Issuer;
export type GcpConfigFormData = GoogleCloudConfigureRequest & Issuer;
