/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model from '@ember-data/model';

export default class MountConfigModel extends Model {
  defaultLeaseTtl: string;
  maxLeaseTtl: string;
  auditNonHmacRequestKeys: string;
  auditNonHmacResponseKeys: string;
  listingVisibility: string;
  passthroughRequestHeaders: string;
  allowedResponseHeaders: string;
  tokenType: string;
  allowedManagedKeys: string;
}
