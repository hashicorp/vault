/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { fromUnixTime, getUnixTime, isSameMonth, isAfter } from 'date-fns';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { formatDateObject } from 'core/utils/client-count-utils';

import type VersionService from 'vault/services/version';
import type ClientsActivityModel from 'vault/models/clients/activity';
import type ClientsConfigModel from 'vault/models/clients/config';
import type StoreService from 'vault/services/store';
import timestamp from 'core/utils/timestamp';

interface Args {
  activity: ClientsActivityModel;
  config: ClientsConfigModel;
  startTimestamp: number;
  endTimestamp: number;
  namespace: string;
  mountPath: string;
  onFilterChange: CallableFunction;
}

export default class ClientsCountsPageComponent extends Component<Args> {
  @service declare readonly version: VersionService;
  @service declare readonly store: StoreService;

  get startTimestampISO() {
    return this.args.startTimestamp ? fromUnixTime(this.args.startTimestamp).toISOString() : null;
  }

  get endTimestampISO() {
    return this.args.endTimestamp ? fromUnixTime(this.args.endTimestamp).toISOString() : null;
  }

  get formattedStartDate() {
    return this.startTimestampISO ? parseAPITimestamp(this.startTimestampISO, 'MMMM yyyy') : null;
  }

  // returns text for empty state message if noActivityData
  get dateRangeMessage() {
    if (this.startTimestampISO && this.endTimestampISO) {
      const endMonth = isSameMonth(
        parseAPITimestamp(this.startTimestampISO) as Date,
        parseAPITimestamp(this.endTimestampISO) as Date
      )
        ? ''
        : `to ${parseAPITimestamp(this.endTimestampISO, 'MMMM yyyy')}`;
      // completes the message 'No data received from { dateRangeMessage }'
      return `from ${parseAPITimestamp(this.startTimestampISO, 'MMMM yyyy')} ${endMonth}`;
    }
    return null;
  }

  get versionText() {
    return this.version.isEnterprise
      ? {
          label: 'Billing start month',
          description:
            'This date comes from your license, and defines when client counting starts. Without this starting point, the data shown is not reliable.',
          title: 'No billing start date found',
          message:
            'In order to get the most from this data, please enter your billing period start month. This will ensure that the resulting data is accurate.',
        }
      : {
          label: 'Client counting start date',
          description:
            'This date is when client counting starts. Without this starting point, the data shown is not reliable.',
          title: 'No start date found',
          message:
            'In order to get the most from this data, please enter a start month above. Vault will calculate new clients starting from that month.',
        };
  }

  get namespaces() {
    return this.args.activity.byNamespace
      ? this.args.activity.byNamespace.map((namespace) => ({
          name: namespace.label,
          id: namespace.label,
        }))
      : [];
  }

  get mountPaths() {
    if (this.namespaces.length) {
      return this.activityForNamespace?.mounts.map((mount) => ({
        id: mount.label,
        name: mount.label,
      }));
    }
    return [];
  }

  get startTimeDiscrepancy() {
    // show banner if startTime returned from activity log (response) is after the queried startTime
    const { activity, config } = this.args;
    const activityStartDateObject = parseAPITimestamp(activity.startTime) as Date;
    const queryStartDateObject = parseAPITimestamp(this.startTimestampISO) as Date;
    const isEnterprise =
      this.startTimestampISO === config.billingStartTimestamp?.toISOString() && this.version.isEnterprise;
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

  get activityForNamespace() {
    const { activity, namespace } = this.args;
    return namespace ? activity.byNamespace.find((ns) => ns.label === namespace) : null;
  }

  get filteredActivity() {
    // return activity counts based on selected namespace and auth mount values
    const { namespace, mountPath, activity } = this.args;
    if (namespace) {
      return mountPath
        ? this.activityForNamespace?.mounts.find((mount) => mount.label === mountPath)
        : this.activityForNamespace;
    }
    return activity.total;
  }

  @action
  onDateChange(dateObject: { dateType: string; monthIdx: string; year: string }) {
    const { dateType, monthIdx, year } = dateObject;
    const { config } = this.args;
    const currentTimestamp = getUnixTime(timestamp.now());

    // converts the selectedDate to unix timestamp for activity query
    const selectedDate = formatDateObject({ monthIdx, year }, dateType === 'endDate');

    if (dateType !== 'cancel') {
      const start_time = {
        reset: getUnixTime(config?.billingStartTimestamp) || null, // clicked 'Current billing period' in calendar widget -> resets to billing start date
        currentMonth: currentTimestamp, // clicked 'Current month' from calendar widget -> defaults to currentTimestamp
        startDate: selectedDate, // from "Edit billing start" modal
      }[dateType];
      // endDate type is selection from calendar widget
      const end_time = dateType === 'endDate' ? selectedDate : currentTimestamp; // defaults to currentTimestamp
      const params = start_time !== undefined ? { start_time, end_time } : { end_time };
      this.args.onFilterChange(params);
    }
  }

  @action
  setFilterValue(type: 'ns' | 'mountPath', [value]: [string | undefined]) {
    const params = { [type]: value };
    // unset mountPath value when namespace is cleared
    if (type === 'ns' && !value) {
      params['mountPath'] = undefined;
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
