/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { Model } from 'vault/app-types';

interface ClientActivityTotals {
  clients: number;
  entity_clients: number;
  non_entity_clients: number;
}

interface ClientActivityNestedCount extends ClientActivityTotals {
  label: string;
}

interface ClientActivityNewClients extends ClientActivityTotals {
  month: string;
  mounts?: ClientActivityNestedCount[];
  namespaces?: ClientActivityNestedCount[];
}

interface ClientActivityNamespace extends ClientActivityNestedCount {
  mounts: ClientActivityNestedCount[];
}

interface ClientActivityNamespaceByKey extends ClientActivityTotals {
  month: 'string';
  mounts_by_key: { [key]: ClientActivityNewClients };
}

interface ClientActivityMonthly extends ClientActivityTotals {
  month: string;
  timestamp: string;
  namespaces: ClientActivityNamespace[];
  namespaces_by_key: { [key: ClientActivityNamespaceByKey] };
  new_clients: ClientActivityNewClients;
}

export default interface ClientsActivityModel extends Model {
  byMonth: ClientActivityMonthly[];
  byNamespace: ClientActivityNamespace[];
  total: ClientActivityTotals;
  startTime: string;
  endTime: string;
  responseTimestamp: string;
}
