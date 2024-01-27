/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { calculateAverage } from 'vault/utils/chart-helpers';

import type { ClientActivityMonthly, ClientActivityTotals } from 'vault/vault/models/clients/activity';

interface Args {
  byMonthActivityData: ClientActivityMonthly[];
  mountPath: string;
  runningTotals: ClientActivityTotals;
  responseTimestamp: string;
}

export default class ClientsTokenMonthlyNewComponent extends Component<Args> {
  runningTotalLegend = [
    { key: 'entity_clients', label: 'entity clients' },
    { key: 'non_entity_clients', label: 'non-entity clients' },
  ];

  get hasAverageNewClients() {
    return (
      typeof this.entityClientData.averageNewClients === 'number' ||
      typeof this.nonEntityClientData.averageNewClients === 'number'
    );
  }

  get byMonthNewClients() {
    if (this.args.byMonthActivityData) {
      return this.args.byMonthActivityData?.map((m) => m.new_clients);
    }
    return null;
  }

  get entityClientData() {
    return {
      runningTotal: this.args.runningTotals.entity_clients,
      averageNewClients: calculateAverage(this.byMonthNewClients, 'entity_clients'),
    };
  }

  get nonEntityClientData() {
    return {
      runningTotal: this.args.runningTotals.non_entity_clients,
      averageNewClients: calculateAverage(this.byMonthNewClients, 'non_entity_clients'),
    };
  }
}
