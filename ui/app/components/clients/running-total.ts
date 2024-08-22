/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import type { ByMonthClients, TotalClients } from 'core/utils/client-count-utils';
import ClientsVersionHistoryModel from 'vault/vault/models/clients/version-history';

interface Args {
  isSecretsSyncActivated: boolean;
  byMonthActivityData: ByMonthClients[];
  isHistoricalMonth: boolean;
  isCurrentMonth: boolean;
  runningTotals: TotalClients;
  upgradesDuringActivity: ClientsVersionHistoryModel[];
  responseTimestamp: string;
  mountPath: string;
}

export default class RunningTotal extends Component<Args> {
  get chartContainerText() {
    const { isSecretsSyncActivated } = this.args;
    return `The total clients in the specified date range, displayed per month. This includes entity, non-entity${
      isSecretsSyncActivated ? ', ACME and secrets sync clients' : ' and ACME clients'
    }. The total client count number is an important consideration for Vault billing.`;
  }

  get runningTotalData() {
    return this.args.byMonthActivityData.map((monthly) => ({
      ...monthly,
      new_clients: monthly.new_clients?.clients,
    }));
  }

  get chartLegend() {
    return [
      { key: 'clients', label: 'total clients' },
      { key: 'new_clients', label: 'new clients' },
    ];
  }
}
