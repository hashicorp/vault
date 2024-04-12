/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { ByMonthClients, MountByKey, NamespaceByKey } from 'core/utils/client-count-utils';
import ActivityComponent from '../activity';

export default class ClientsAcmePageComponent extends ActivityComponent {
  title = 'ACME usage';
  get description() {
    return `This data can be used to understand how many ACME clients have been used for the queried ${
      this.isDateRange ? 'date range' : 'month'
    }. Each ACME request is counted as one client.`;
  }

  get hasMonthlyData() {
    return this.byMonthActivityData.any((month: ByMonthClients | NamespaceByKey | MountByKey | undefined) => {
      if (!month) return false;
      return !!month.acme_clients;
    });
  }
}
