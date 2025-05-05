/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { capitalize } from '@ember/string';

import type { DestinationType } from 'vault/sync';

type TypeKey = 'AwsSm' | 'AzureKv' | 'GcpSm' | 'Gh' | 'VercelProject';

// unfortunately type is not a param for sync destination api methods
// it seems that for documentation purposes it was intentionally set up this way
export default function apiMethodResolver(action: 'read' | 'write' | 'delete', type: DestinationType) {
  const formattedType = type.split('-').reduce((str, word) => str.concat(capitalize(word)), '') as TypeKey;
  const method = `system${capitalize(action)}SyncDestinations${formattedType}Name`;

  return method as `system${Capitalize<typeof action>}SyncDestinations${typeof formattedType}Name`;
}
