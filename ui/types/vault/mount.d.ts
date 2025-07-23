/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export type MountConfig = {
  force_no_cache?: boolean;
  listing_visibility?: string | boolean;
  default_lease_ttl?: number;
  max_lease_ttl?: number;
  allowed_managed_keys?: string[];
  audit_non_hmac_request_keys?: string[];
  audit_non_hmac_response_keys?: string[];
  passthrough_request_headers?: string[];
  allowed_response_headers?: string[];
  identity_token_key?: string;
};

export type MountOptions = {
  version: number;
};

export type Mount = {
  path: string;
  accessor: string;
  config: MountConfig;
  description: string;
  external_entropy_access: boolean;
  local: boolean;
  options?: MountOptions;
  plugin_version: string;
  running_plugin_version: string;
  running_sha256: string;
  seal_wrap: boolean;
  type: string;
  uuid: string;
};
