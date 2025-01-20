/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// base component for counts child routes that can be extended as needed
// contains getters that filter and extract data from activity model for use in charts

import Component from '@glimmer/component';
import { isAfter, isBefore, isSameMonth, fromUnixTime } from 'date-fns';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { calculateAverage } from 'vault/utils/chart-helpers';

import type ClientsActivityModel from 'vault/models/clients/activity';
import type {
  ClientActivityNewClients,
  ClientActivityMonthly,
  ClientActivityResourceByKey,
} from 'vault/models/clients/activity';
import type ClientsVersionHistoryModel from 'vault/models/clients/version-history';

interface Args {
  activity: ClientsActivityModel;
  versionHistory: ClientsVersionHistoryModel[];
  startTimestamp: number;
  endTimestamp: number;
  namespace: string;
  mountPath: string;
}

export default class ClientsActivityComponent extends Component<Args> {
  average = (
    data:
      | ClientActivityMonthly[]
      | (ClientActivityResourceByKey | undefined)[]
      | (ClientActivityNewClients | undefined)[]
      | undefined,
    key: string
  ) => {
    return calculateAverage(data, key);
  };

  get startTimeISO() {
    return fromUnixTime(this.args.startTimestamp).toISOString();
  }

  get endTimeISO() {
    return fromUnixTime(this.args.endTimestamp).toISOString();
  }

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
      .map((m) => m.namespaces_by_key[namespace as keyof typeof m.namespaces_by_key])
      .filter((d) => d !== undefined);

    if (!mountPath) {
      return namespaceData.length === 0 ? undefined : namespaceData;
    }

    const mountData = mountPath
      ? namespaceData.map((namespace) => namespace?.mounts_by_key[mountPath]).filter((d) => d !== undefined)
      : namespaceData;

    return mountData.length === 0 ? undefined : mountData;
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

  get upgradeDuringActivity() {
    const { versionHistory, activity } = this.args;
    if (versionHistory) {
      // filter for upgrade data of noteworthy upgrades (1.9 and/or 1.10)
      const upgradeVersionHistory = versionHistory.filter(
        ({ version }) => version.match('1.9') || version.match('1.10')
      );
      if (upgradeVersionHistory.length) {
        const activityStart = parseAPITimestamp(activity.startTime) as Date;
        const activityEnd = parseAPITimestamp(activity.endTime) as Date;
        // filter and return all upgrades that happened within date range of queried activity
        const upgradesWithinData = upgradeVersionHistory.filter(({ timestampInstalled }) => {
          const upgradeDate = parseAPITimestamp(timestampInstalled) as Date;
          return isAfter(upgradeDate, activityStart) && isBefore(upgradeDate, activityEnd);
        });
        return upgradesWithinData.length === 0 ? null : upgradesWithinData;
      }
    }
    return null;
  }

  // (object) single month new client data with total counts + array of namespace breakdown
  get newClientCounts() {
    if (this.isDateRange || !this.byMonthActivityData) {
      return null;
    }
    return this.byMonthActivityData[0]?.new_clients;
  }

  // total client data for horizontal bar chart in attribution component
  get totalClientAttribution() {
    const { namespace, activity } = this.args;
    if (namespace) {
      return this.filteredActivityByNamespace?.mounts || null;
    } else {
      return activity.byNamespace || null;
    }
  }

  // new client data for horizontal bar chart
  get newClientAttribution() {
    // new client attribution only available in a single, historical month (not a date range or current month)
    if (this.isDateRange || this.isCurrentMonth) return null;

    if (this.args.namespace) {
      return this.newClientCounts?.mounts || null;
    } else {
      return this.newClientCounts?.namespaces || null;
    }
  }

  get hasAttributionData() {
    const { mountPath, namespace } = this.args;
    if (!mountPath) {
      if (namespace) {
        const mounts = this.filteredActivityByNamespace?.mounts?.map((mount) => ({
          id: mount.label,
          name: mount.label,
        }));
        return mounts && mounts.length > 0;
      }
      return !!this.totalClientAttribution && this.totalUsageCounts && this.totalUsageCounts.clients !== 0;
    }

    return false;
  }

  get upgradeExplanation() {
    if (this.upgradeDuringActivity) {
      if (this.upgradeDuringActivity.length === 1) {
        const version = this.upgradeDuringActivity[0]?.version || '';
        if (version.match('1.9')) {
          return ' How we count clients changed in 1.9, so keep that in mind when looking at the data.';
        }
        if (version.match('1.10')) {
          return ' We added monthly breakdowns and mount level attribution starting in 1.10, so keep that in mind when looking at the data.';
        }
      }
      // return combined explanation if spans multiple upgrades
      return ' How we count clients changed in 1.9 and we added monthly breakdowns and mount level attribution starting in 1.10. Keep this in mind when looking at the data.';
    }
    return null;
  }
}
