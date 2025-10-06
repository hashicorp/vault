/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export type MountConfig = {
  allowed_managed_keys?: string[];
  allowed_response_headers?: string[];
  audit_non_hmac_request_keys?: string[];
  audit_non_hmac_response_keys?: string[];
  default_lease_ttl?: number | string;
  force_no_cache?: boolean;
  identity_token_key?: string;
  listing_visibility?: string | boolean;
  max_lease_ttl?: number | string;
  passthrough_request_headers?: string[];
  token_type?: string;
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
