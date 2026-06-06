/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { DestinationType } from 'sync/utils/constants';
import { DestinationName, DestinationRoleTypeOption } from 'vault/sync';

export interface SyncDestination {
  name: DestinationName;
  type: DestinationType;
  icon: 'aws-color' | 'azure-color' | 'gcp-color' | 'github-color' | 'vercel-color';
  category: 'cloud' | 'dev-tools';
  maskedParams: Array<string>;
  readonlyParams: Array<string>;
  defaultValues: object;
  roleTypeOptions?: Array<DestinationRoleTypeOption>;
}
