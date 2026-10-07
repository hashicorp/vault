/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { baseResourceFactory } from 'vault/resources/base-factory';

type OidcKeyData = {
  name: string;
  algorithm: string;
  rotation_period: number | string;
  verification_ttl: number | string;
  allowed_client_ids?: string[];
};

/** Read-only data for an OIDC key detail view. */
export default class OidcKeyResource extends baseResourceFactory<OidcKeyData>() {}
