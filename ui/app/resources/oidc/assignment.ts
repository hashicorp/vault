/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { baseResourceFactory } from 'vault/resources/base-factory';

type OidcAssignmentData = {
  name: string;
  entity_ids?: string[];
  group_ids?: string[];
  entities?: string[];
  groups?: string[];
};

/** Read-only data for an OIDC assignment detail view. */
export default class OidcAssignmentResource extends baseResourceFactory<OidcAssignmentData>() {}
