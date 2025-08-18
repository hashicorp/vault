/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

import type { ByMonthNewClients, TotalClients } from 'core/utils/client-count-utils';
import type ClientsVersionHistoryModel from 'vault/vault/models/clients/version-history';

interface Args {
  isSecretsSyncActivated: boolean;
  byMonthNewClients: ByMonthNewClients[];
  isHistoricalMonth: boolean;
  isCurrentMonth: boolean;
  runningTotals: TotalClients;
  upgradesDuringActivity: ClientsVersionHistoryModel[];
  responseTimestamp: string;
  mountPath: string;
}

export default class RunningTotal extends Component<Args> {
  @tracked showStacked = false;

  get chartContainerText() {
    const { isSecretsSyncActivated } = this.args;
    return `The total clients in the specified date range, displayed per month. This includes entity, non-entity${
      isSecretsSyncActivated ? ', ACME and secrets sync clients' : ' and ACME clients'
    }. The total client count number is an important consideration for Vault billing.`;
  }

  get runningTotalData() {
    return this.args.byMonthNewClients.map((monthly) => ({
      ...monthly,
      new_clients: monthly.clients,
    }));
  }

  get chartLegend() {
    if (this.showStacked) {
      return [
        { key: 'entity_clients', label: 'entity clients' },
        { key: 'non_entity_clients', label: 'non-entity clients' },
        ...(this.args.isSecretsSyncActivated ? [{ key: 'secret_syncs', label: 'secret sync clients' }] : []),
        { key: 'acme_clients', label: 'acme clients' },
      ];
    }
    return [{ key: 'new_clients', label: 'new clients' }];
  }
}
