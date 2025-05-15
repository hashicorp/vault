/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { DestinationName, DestinationType } from 'vault/sync';

export interface SyncDestination {
  name: DestinationName;
  type: DestinationType;
  icon: 'aws-color' | 'azure-color' | 'gcp-color' | 'github-color' | 'vercel-color';
  category: 'cloud' | 'dev-tools';
  maskedParams: Array<string>;
  readonlyParams: Array<string>;
  defaultValues: object;
}
