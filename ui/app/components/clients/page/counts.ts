/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { filterVersionHistory } from 'core/utils/client-counts/helpers';

import type AdapterError from '@ember-data/adapter/error';
import type FlagsService from 'vault/services/flags';
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

  get error() {
    const { httpStatus, message, path } = this.args.activityError || {};
    let title = 'Error',
      text = message;

    if (httpStatus === 403) {
      const endpoint = path ? `the ${path} endpoint` : 'this endpoint';
      title = 'You are not authorized';
      text = `You must be granted permissions to view this page. Ask your administrator if you think you should have access to ${endpoint}.`;
    }

    return { title, text, httpStatus };
  }

  get formattedStartDate() {
    return this.args.startTimestamp ? parseAPITimestamp(this.args.startTimestamp, 'MMMM yyyy') : null;
  }

  get formattedEndDate() {
    return this.args.endTimestamp ? parseAPITimestamp(this.args.endTimestamp, 'MMMM yyyy') : null;
  }

  get formattedBillingStartDate() {
    if (this.args.config?.billingStartTimestamp) {
      return this.args.config.billingStartTimestamp.toISOString();
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

  // the dashboard should show sync tab if the flag is on or there's data
  get hasSecretsSyncClients(): boolean {
    return this.args.activity?.total?.secret_syncs > 0;
  }

  @action
  onDateChange(params: { start_time: string; end_time: string }) {
    this.args.onFilterChange(params);
  }
}
