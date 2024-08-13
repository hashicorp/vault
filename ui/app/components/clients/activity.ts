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
  ByMonthNewClients,
  MountNewClients,
  NamespaceByKey,
  NamespaceNewClients,
} from 'core/utils/client-count-utils';

interface Args {
  activity: ClientsActivityModel;
  versionHistory: ClientsVersionHistoryModel[];
  startTimestamp: string;
  endTimestamp: string;
  namespace: string;
  mountPath: string;
}

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
    const { activity, namespace } = this.args;
    return namespace ? this.filteredActivityByMonth : activity.byMonth;
  }

  get byMonthNewClients() {
    return this.byMonthActivityData ? this.byMonthActivityData?.map((m) => m?.new_clients) : [];
  }

  get filteredActivityByMonth() {
    const { namespace, mountPath, activity } = this.args;
    if (!namespace && !mountPath) {
      return activity.byMonth;
    }
    const namespaceData = activity.byMonth
      ?.map((m) => m.namespaces_by_key[namespace])
      .filter((d) => d !== undefined);

    if (!mountPath) {
      return namespaceData || [];
    }

    const mountData = namespaceData
      ?.map((namespace) => namespace?.mounts_by_key[mountPath])
      .filter((d) => d !== undefined);

    return mountData || [];
  }

  get filteredActivityByNamespace() {
    const { namespace, activity } = this.args;
    return activity.byNamespace.find((ns) => ns.label === namespace);
  }

  get filteredActivityByAuthMount() {
    return this.filteredActivityByNamespace?.mounts?.find((mount) => mount.label === this.args.mountPath);
  }

  get filteredActivity() {
    return this.args.mountPath ? this.filteredActivityByAuthMount : this.filteredActivityByNamespace;
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

  // (object) top level TOTAL client counts for given date range
  get totalUsageCounts() {
    const { namespace, activity } = this.args;
    return namespace ? this.filteredActivity : activity.total;
  }

  get upgradesDuringActivity() {
    const { versionHistory, activity } = this.args;
    return filterVersionHistory(versionHistory, activity.startTime, activity.endTime);
  }
}
