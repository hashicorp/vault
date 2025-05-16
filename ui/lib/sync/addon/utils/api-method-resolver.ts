/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Util that resolves the API client method name for a given action and destination type
 * Typically both type and name would be defined as params -> /sys/sync/destinations/{type}/{name}
 * This would result in a method on the API client like systemReadSyncDestinations('aws-sm', 'my-destination')
 * This unfortunately is not the case and type was hardcoded in the path -> /sys/sync/destinations/aws-sm/{name}
 * This results in GET, POST and DELETE methods for each destination type -> systemReadSyncDestinationsAwsSmName('my-destination')
 * Since these are used in numerous places, this util was created to more easily resolve the method name
 */

import { capitalize, classify } from '@ember/string';

import type { DestinationType } from 'vault/sync';

type TypeKey = 'AwsSm' | 'AzureKv' | 'GcpSm' | 'Gh' | 'VercelProject';

export default function apiMethodResolver(
  action: 'read' | 'write' | 'delete' | 'patch',
  type: DestinationType
) {
  const method = `system${capitalize(action)}SyncDestinations${classify(type)}Name`;
  return method as `system${Capitalize<typeof action>}SyncDestinations${TypeKey}Name`;
}
