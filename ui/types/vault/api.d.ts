/**
 * Copyright IBM Corp. 2016, 2026
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

export interface ApiBaseErrorResponse {
  data?: {
    error?: string;
  };
  error?: string;
  errors?: string[];
  message?: string;
  [key: string]: unknown;
}

export interface ControlGroupErrorResponse extends ApiBaseErrorResponse, WrapInfo {
  isControlGroupError: true;
}

export interface ApiStandardErrorResponse extends ApiBaseErrorResponse {
  isControlGroupError?: false;
}

export type ApiErrorResponse = ControlGroupErrorResponse | ApiStandardErrorResponse;

export interface ApiParsedError {
  message: string;
  path?: string;
  response?: ApiErrorResponse;
  status?: number;
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
    }
  | {
      recoverSnapshotId: string;
    }
  | {
      recoverSourcePath: string;
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
    }
  | {
      'X-Vault-Snapshot-Id': string;
    }
  | {
      'X-Vault-Recover-Source-Path': string;
    };
