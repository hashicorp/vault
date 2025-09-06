/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { isSameMonth, isAfter } from 'date-fns';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { filterVersionHistory } from 'core/utils/client-count-utils';

import type AdapterError from '@ember-data/adapter/error';
import type FlagsService from 'vault/services/flags';
import type Store from '@ember-data/store';
import type VersionService from 'vault/services/version';
import type ClientsActivityModel from 'vault/models/clients/activity';
import type ClientsConfigModel from 'vault/models/clients/config';
import type ClientsVersionHistoryModel from 'vault/models/clients/version-history';

interface Args {
  activity: ClientsActivityModel;
  activityError?: AdapterError;
  config: ClientsConfigModel;
  endTimestamp: string; // ISO format
  onFilterChange: CallableFunction;
  startTimestamp: string; // ISO format
  versionHistory: ClientsVersionHistoryModel[];
}

export default class ClientsCountsPageComponent extends Component<Args> {
  @service declare readonly flags: FlagsService;
  @service declare readonly version: VersionService;
  @service declare readonly store: Store;

  get formattedStartDate() {
    return this.args.startTimestamp ? parseAPITimestamp(this.args.startTimestamp, 'MMMM yyyy') : null;
  }

  get formattedEndDate() {
    return this.args.endTimestamp ? parseAPITimestamp(this.args.endTimestamp, 'MMMM yyyy') : null;
  }

  get formattedBillingStartDate() {
    return this.args.config.billingStartTimestamp.toISOString();
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
    return filterVersionHistory(versionHistory, activity?.startTime, activity?.endTime);
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
}
