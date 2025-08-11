/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// base component for counts child routes that can be extended as needed
// contains getters that filter and extract data from activity model for use in charts

import Component from '@glimmer/component';
import { cached } from '@glimmer/tracking';
import { isSameMonth } from 'date-fns';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import {
  filterByMonthDataForMount,
  filteredTotalForMount,
  filterVersionHistory,
} from 'core/utils/client-count-utils';
import { service } from '@ember/service';
import { sanitizePath } from 'core/utils/sanitize-path';

import type ClientsActivityModel from 'vault/models/clients/activity';
import type ClientsVersionHistoryModel from 'vault/models/clients/version-history';
import type { MountClients, TotalClients } from 'core/utils/client-count-utils';
import type NamespaceService from 'vault/services/namespace';

interface Args {
  activity: ClientsActivityModel;
  versionHistory: ClientsVersionHistoryModel[];
  startTimestamp: string;
  endTimestamp: string;
  namespace: string;
  mountPath: string;
  mountType: string;
  onFilterChange: CallableFunction;
}

export default class ClientsActivityComponent extends Component<Args> {
  @service declare readonly namespace: NamespaceService;

  @cached
  get namespaceLabels() {
    return this.args.activity.byNamespace.map((n) => n.label);
  }

  @cached
  get mounts() {
    return this.args.activity.byNamespace.map((n) => n.mounts).flat();
  }

  @cached
  get mountPaths() {
    return [...new Set(this.mounts.map((m: MountClients) => m.label))];
  }

  @cached
  get mountTypes() {
    return [...new Set(this.mounts.map((m: MountClients) => m.mount_type))];
  }

  // path of the filtered namespace OR current one, for filtering relevant data
  get namespacePathForFilter() {
    const { namespace } = this.args;
    const currentNs = this.namespace.currentNamespace;
    return sanitizePath(namespace || currentNs || 'root');
  }

  get byMonthNewClients() {
    const { activity, mountPath } = this.args;
    const nsPath = this.namespacePathForFilter;

    const data = mountPath
      ? filterByMonthDataForMount(activity.byMonth, nsPath, mountPath)
      : activity.byMonth;

    return data ? data?.map((m) => m?.new_clients) : [];
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
  get totalUsageCounts(): TotalClients {
    const { namespace, activity, mountPath } = this.args;
    // only do this if we have a mountPath filter.
    // namespace is filtered on API layer
    if (activity?.byNamespace && namespace && mountPath) {
      return filteredTotalForMount(activity.byNamespace, namespace, mountPath);
    }
    return activity?.total;
  }

  get upgradesDuringActivity() {
    const { versionHistory, activity } = this.args;
    return filterVersionHistory(versionHistory, activity.startTime, activity.endTime);
  }
}
