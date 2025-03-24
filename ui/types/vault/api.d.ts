/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { ErrorContext } from '@hashicorp/vault-client-typescript';

// re-exporting for convenience since it is associated to ApiError
export { ErrorContext };
export interface ApiError {
  httpStatus: number;
  path: string;
  message: string;
  errors: Array<string | { [key: string]: unknown; title?: string; message?: string }>;
  data?: {
    [key: string]: unknown;
    error?: string;
  };
}

export interface WrapInfo {
  accessor: string;
  creation_path: string;
  creation_time: string;
  wrapped_accessor: string;
  token: string;
  ttl: number;
}

export interface ApiResponse {
  auth: unknown;
  data: unknown;
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
