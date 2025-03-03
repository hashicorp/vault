/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { isSameMonth, isAfter } from 'date-fns';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { filteredTotalForMount, filterVersionHistory, TotalClients } from 'core/utils/client-count-utils';
import { sanitizePath } from 'core/utils/sanitize-path';

import type AdapterError from '@ember-data/adapter/error';
import type FlagsService from 'vault/services/flags';
import type Store from '@ember-data/store';
import type VersionService from 'vault/services/version';
import type ClientsActivityModel from 'vault/models/clients/activity';
import type ClientsConfigModel from 'vault/models/clients/config';
import type ClientsVersionHistoryModel from 'vault/models/clients/version-history';
import type NamespaceService from 'vault/services/namespace';

interface Args {
  activity: ClientsActivityModel;
  activityError?: AdapterError;
  config: ClientsConfigModel;
  endTimestamp: string; // ISO format
  mountPath: string;
  namespace: string;
  onFilterChange: CallableFunction;
  startTimestamp: string; // ISO format
  versionHistory: ClientsVersionHistoryModel[];
}

export default class ClientsCountsPageComponent extends Component<Args> {
  @service declare readonly flags: FlagsService;
  @service declare readonly version: VersionService;
  @service declare readonly namespace: NamespaceService;
  @service declare readonly store: Store;

  get formattedStartDate() {
    return this.args.startTimestamp ? parseAPITimestamp(this.args.startTimestamp, 'MMMM yyyy') : null;
  }

  // returns text for empty state message if noActivityData
  get dateRangeMessage() {
    if (this.args.startTimestamp && this.args.endTimestamp) {
      const endMonth = isSameMonth(
        parseAPITimestamp(this.args.startTimestamp) as Date,
        parseAPITimestamp(this.args.endTimestamp) as Date
      )
        ? ''
        : `to ${parseAPITimestamp(this.args.endTimestamp, 'MMMM yyyy')}`;
      // completes the message 'No data received from { dateRangeMessage }'
      return `from ${parseAPITimestamp(this.args.startTimestamp, 'MMMM yyyy')} ${endMonth}`;
    }
    return null;
  }

  // passed into page-header for the export modal alert
  get upgradesDuringActivity() {
    const { versionHistory, activity } = this.args;
    return filterVersionHistory(versionHistory, activity.startTime, activity.endTime);
  }

  get upgradeExplanations() {
    if (this.upgradesDuringActivity.length) {
      return this.upgradesDuringActivity.map((upgrade: ClientsVersionHistoryModel) => {
        let explanation;
        const date = parseAPITimestamp(upgrade.timestampInstalled, 'MMM d, yyyy');
        const version = upgrade.version || '';
        switch (true) {
          case version.includes('1.9'):
            explanation =
              '- We introduced changes to non-entity token and local auth mount logic for client counting in 1.9.';
            break;
          case version.includes('1.10'):
            explanation = '- We added monthly breakdowns and mount level attribution starting in 1.10.';
            break;
          case version.includes('1.17'):
            explanation = '- We separated ACME clients from non-entity clients starting in 1.17.';
            break;
          default:
            explanation = '';
            break;
        }
        return `${version} (upgraded on ${date}) ${explanation}`;
      });
    }
    return null;
  }

  get versionText() {
    return this.version.isEnterprise
      ? {
          title: 'No billing start date found',
          message:
            'In order to get the most from this data, please enter your billing period start month. This will ensure that the resulting data is accurate.',
        }
      : {
          title: 'No start date found',
          message:
            'In order to get the most from this data, please enter a start month above. Vault will calculate new clients starting from that month.',
        };
  }

  // path of the filtered namespace OR current one, for filtering relevant data
  get namespacePathForFilter() {
    const { namespace } = this.args;
    const currentNs = this.namespace.currentNamespace;
    return sanitizePath(namespace || currentNs || 'root');
  }

  // activityForNamespace gets the byNamespace data for the selected or current namespace so we can get the list of mounts from that namespace for attribution
  get activityForNamespace() {
    const { activity } = this.args;
    const nsPath = this.namespacePathForFilter;
    // we always return activity for namespace, either the selected filter or the current
    return activity?.byNamespace?.find((ns) => sanitizePath(ns.label) === nsPath);
  }

  // duplicate of the method found in the activity component, so that we render the child only when there is activity to view
  get totalUsageCounts(): TotalClients {
    const { namespace, mountPath, activity } = this.args;
    if (mountPath) {
      // only do this if we have a mountPath filter.
      // namespace is filtered on API layer
      return filteredTotalForMount(activity.byNamespace, namespace, mountPath);
    }
    return activity?.total;
  }

  // namespace list for the search-select filter
  get namespaces() {
    return this.args.activity?.byNamespace
      ? this.args.activity.byNamespace
          .map((namespace) => ({
            name: namespace.label,
            id: namespace.label,
          }))
          .filter((ns) => sanitizePath(ns.name) !== this.namespacePathForFilter)
      : [];
  }

  // mounts within the current/filtered namespace for the sesarch-select filter
  get mountPaths() {
    return (
      this.activityForNamespace?.mounts.map((mount) => ({
        id: mount.label,
        name: mount.label,
      })) || []
    );
  }

  // banner contents shown if startTime returned from activity API (which matches the first month with data) is after the queried startTime
  get startTimeDiscrepancy() {
    const { activity, config } = this.args;
    const activityStartDateObject = parseAPITimestamp(activity.startTime) as Date;
    const queryStartDateObject = parseAPITimestamp(this.args.startTimestamp) as Date;
    const isEnterprise =
      this.args.startTimestamp === config.billingStartTimestamp?.toISOString() && this.version.isEnterprise;
    const message = isEnterprise ? 'Your license start date is' : 'You requested data from';

    if (
      isAfter(activityStartDateObject, queryStartDateObject) &&
      !isSameMonth(activityStartDateObject, queryStartDateObject)
    ) {
      return `${message} ${this.formattedStartDate}.
        We only have data from ${parseAPITimestamp(activity.startTime, 'MMMM yyyy')},
        and that is what is being shown here.`;
    } else {
      return null;
    }
  }

  // the dashboard should show sync tab if the flag is on or there's data
  get hasSecretsSyncClients(): boolean {
    return this.args.activity?.total?.secret_syncs > 0;
  }

  @action
  onDateChange(params: { start_time: number | undefined; end_time: number | undefined }) {
    this.args.onFilterChange(params);
  }

  @action
  setFilterValue(type: 'ns' | 'mountPath', [value]: [string | undefined]) {
    const params = { [type]: value };
    if (type === 'ns' && !value) {
      // unset mountPath value when namespace is cleared
      params['mountPath'] = undefined;
    } else if (type === 'mountPath' && !this.args.namespace) {
      // set namespace when mountPath set without namespace already set
      params['ns'] = this.namespacePathForFilter;
    }
    this.args.onFilterChange(params);
  }

  @action resetFilters() {
    this.args.onFilterChange({
      start_time: undefined,
      end_time: undefined,
      ns: undefined,
      mountPath: undefined,
    });
  }
}
