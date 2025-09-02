/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { Model } from 'vault/app-types';

export default interface ClientsVersionHistoryModel extends Model {
  version: string;
  previousVersion: string;
  timestampInstalled: string;
}
