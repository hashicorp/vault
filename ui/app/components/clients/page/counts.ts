/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { fromUnixTime, isSameMonth, isAfter } from 'date-fns';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { formatDateObject } from 'core/utils/client-count-utils';

import type VersionService from 'vault/services/version';
import type ClientsActivityModel from 'vault/models/clients/activity';
import type ClientsConfigModel from 'vault/models/clients/config';
import type RouterService from '@ember/routing/router-service';
import type StoreService from 'vault/services/store';

interface Args {
  activity: ClientsActivityModel;
  config: ClientsConfigModel;
  startTimestamp: number;
  endTimestamp: number;
  currentTimestamp: number;
  namespace: string;
  authMount: string;
}

export default class ClientFiltersHeader extends Component<Args> {
  @service declare readonly version: VersionService;
  @service declare readonly router: RouterService;
  @service declare readonly store: StoreService;

  get startDate() {
    return this.args.startTimestamp ? fromUnixTime(this.args.startTimestamp).toISOString() : null;
  }

  get endDate() {
    return this.args.endTimestamp ? fromUnixTime(this.args.endTimestamp).toISOString() : null;
  }

  get formattedStartDate() {
    return this.startDate ? parseAPITimestamp(this.startDate, 'MMMM yyyy') : null;
  }

  // returns text for empty state message if noActivityData
  get dateRangeMessage() {
    if (this.startDate && this.endDate) {
      const endMonth = isSameMonth(
        parseAPITimestamp(this.startDate) as Date,
        parseAPITimestamp(this.endDate) as Date
      )
        ? ''
        : `to ${parseAPITimestamp(this.endDate, 'MMMM yyyy')}`;
      // completes the message 'No data received from { dateRangeMessage }'
      return `from ${parseAPITimestamp(this.startDate, 'MMMM yyyy')} ${endMonth}`;
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

  get authMounts() {
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
    const queryStartDateObject = parseAPITimestamp(this.startDate) as Date;
    const isEnterprise =
      this.startDate === config.billingStartTimestamp.toISOString() && this.version.isEnterprise;
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
    const { namespace, authMount, activity } = this.args;
    if (namespace) {
      return authMount
        ? this.activityForNamespace?.mounts.find((mount) => mount.label === authMount)
        : this.activityForNamespace;
    }
    return activity.total;
  }

  @action
  onDateChange(dateObject: { dateType: string; monthIdx: string; year: string }) {
    const { dateType, monthIdx, year } = dateObject;
    const { currentTimestamp } = this.args;
    const selectedDate = formatDateObject({ monthIdx, year }, dateType === 'endDate');

    if (dateType !== 'cancel') {
      const start_time = {
        reset: null, // clicked 'Current billing period' in calendar widget -> defaults to licenseStartTime
        currentMonth: currentTimestamp, // clicked 'Current month' from calendar widget
        startDate: selectedDate, // from "Edit billing start" modal
      }[dateType];
      // endDate was selections from calendar widget
      const end_time = dateType === 'endDate' ? selectedDate : currentTimestamp;
      const queryParams = start_time !== undefined ? { start_time, end_time } : { end_time };
      this.router.transitionTo({ queryParams });
    }
  }

  @action
  setFilterValue(type: 'ns' | 'authMount', [value]: [string]) {
    this.router.transitionTo({ queryParams: { [type]: value } });
  }

  @action resetFilters() {
    this.router.transitionTo({
      queryParams: { start_time: null, end_time: null, ns: null, authMount: null },
    });
  }
}
