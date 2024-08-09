/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// base component for counts child routes that can be extended as needed
// contains getters that filter and extract data from activity model for use in charts

import Component from '@glimmer/component';
import { isSameMonth } from 'date-fns';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { calculateAverage } from 'vault/utils/chart-helpers';
import { filterVersionHistory } from 'core/utils/client-count-utils';

import type ClientsActivityModel from 'vault/models/clients/activity';
import type ClientsVersionHistoryModel from 'vault/models/clients/version-history';
import type {
  ByMonthClients,
  ByMonthNewClients,
  MountByKey,
  MountNewClients,
  NamespaceByKey,
  NamespaceNewClients,
  TotalClients,
} from 'core/utils/client-count-utils';

interface Args {
  activity: ClientsActivityModel;
  versionHistory: ClientsVersionHistoryModel[];
  startTimestamp: string;
  endTimestamp: string;
  namespace: string;
  mountPath: string;
  filteredByMonth: ByMonthClients[] | MountByKey[];
  filteredTotals: TotalClients | undefined;
}

// Component class extended by Clients::Page::* components
export default class ClientsActivityComponent extends Component<Args> {
  average = (
    data:
      | (ByMonthNewClients | NamespaceNewClients | MountNewClients | undefined)[]
      | (NamespaceByKey | undefined)[],
    key: string
  ) => {
    return calculateAverage(data, key);
  };

  get byMonthActivityData() {
    return this.args.filteredByMonth;
  }

  get byMonthNewClients() {
    return this.byMonthActivityData ? this.byMonthActivityData?.map((m) => m?.new_clients) : [];
  }

  get isCurrentMonth() {
    const { activity } = this.args;
    const current = parseAPITimestamp(activity.responseTimestamp) as Date;
    const start = parseAPITimestamp(activity.startTime) as Date;
    const end = parseAPITimestamp(activity.endTime) as Date;
    return isSameMonth(start, current) && isSameMonth(end, current);
  }

  get isDateRange() {
    const { activity } = this.args;
    return !isSameMonth(
      parseAPITimestamp(activity.startTime) as Date,
      parseAPITimestamp(activity.endTime) as Date
    );
  }

  // (object) TOTAL client counts for given date range (filtered)
  get totalUsageCounts() {
    return this.args.filteredTotals;
  }

  get upgradesDuringActivity() {
    const { versionHistory, activity } = this.args;
    return filterVersionHistory(versionHistory, activity.startTime, activity.endTime);
  }

  // (object) single month new client data with total counts and array of
  // either namespaces or mounts
  get newClientCounts() {
    if (this.isDateRange || this.byMonthActivityData.length === 0) {
      return null;
    }

    return this.byMonthActivityData[0]?.new_clients;
  }
}
