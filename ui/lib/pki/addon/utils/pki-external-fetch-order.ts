/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import type ApiService from 'vault/services/api';
import type { ApiParsedError } from 'vault/api';
import { PkiExternalCaReadRoleOrderFetchCertResponse } from '@hashicorp/vault-client-typescript';

/**
 * Fetches a certificate for a role order. Errors are captured rather than thrown
 * so the caller can display why this request failed without hiding details from
 * successful requests.
 */
export async function fetchRoleOrderCert(
  api: ApiService,
  roleName: string,
  orderId: string,
  mountPath: string
) {
  let details: PkiExternalCaReadRoleOrderFetchCertResponse | undefined;
  let error: ApiParsedError | undefined;
  try {
    details = await api.secrets.pkiExternalCaReadRoleOrderFetchCert(roleName, orderId, mountPath);
  } catch (e) {
    error = await api.parseError(e);
  }
  return { details, error };
}
