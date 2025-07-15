/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export interface WrapInfo {
  accessor: string;
  creation_path: string;
  creation_time: string;
  wrapped_accessor: string;
  token: string;
  ttl: number;
}

export interface ApiResponse {
  data: Record<string, unknown> | null;
  auth: Record<string, unknown> | null;
  lease_duration: number;
  lease_id: string;
  mount_type: string;
  renewable: boolean;
  request_id: string;
  warnings: Array<string> | null;
  wrap_info: WrapInfo | null;
}

export type HeaderMap =
  | {
      namespace: string;
    }
  | {
      token: string;
    }
  | {
      wrap: string;
    };

export type XVaultHeaders =
  | {
      'X-Vault-Namespace': string;
    }
  | {
      'X-Vault-Token': string;
    }
  | {
      'X-Vault-Wrap-TTL': string;
    };
