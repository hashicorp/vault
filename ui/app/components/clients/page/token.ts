/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ActivityComponent from '../activity';

import type {
  ByMonthNewClients,
  MountNewClients,
  NamespaceByKey,
  NamespaceNewClients,
} from 'core/utils/client-count-utils';

export default class ClientsTokenPageComponent extends ActivityComponent {
  legend = [
    { key: 'entity_clients', label: 'entity clients' },
    { key: 'non_entity_clients', label: 'non-entity clients' },
  ];

  calculateClientAverages(
    dataset:
      | (NamespaceByKey | undefined)[]
      | (ByMonthNewClients | NamespaceNewClients | MountNewClients | undefined)[]
  ) {
    return ['entity_clients', 'non_entity_clients'].reduce((count, key) => {
      const average = this.average(dataset, key);
      return (count += average || 0);
    }, 0);
  }

  get averageTotalClients() {
    return this.calculateClientAverages(this.byMonthActivityData);
  }

  get averageNewClients() {
    return this.calculateClientAverages(this.byMonthNewClients);
  }

  get tokenStats() {
    if (this.totalUsageCounts) {
      const { entity_clients, non_entity_clients } = this.totalUsageCounts;
      return {
        total: entity_clients + non_entity_clients,
        entity_clients,
        non_entity_clients,
      };
    }
    return null;
  }
}
