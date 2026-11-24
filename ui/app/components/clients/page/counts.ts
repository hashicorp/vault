/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { filterVersionHistory } from 'core/utils/client-counts/helpers';

import type FlagsService from 'vault/services/flags';
import type VersionService from 'vault/services/version';
import type { VersionHistory } from 'vault/client-counts';
import type { Activity } from 'vault/client-counts/activity-api';
import type { InternalClientActivityReadConfigurationResponse } from '@hashicorp/vault-client-typescript';

interface Args {
  activity: Activity;
  config: InternalClientActivityReadConfigurationResponse;
  canUpdateConfig: boolean;
  endTimestamp: Date;
  onFilterChange: CallableFunction;
  startTimestamp: Date;
  versionHistory: VersionHistory[];
  responseTimestamp: Date;
}

export default class ClientsCountsPageComponent extends Component<Args> {
  @service declare readonly flags: FlagsService;
  @service declare readonly version: VersionService;

  get trackingDisabled() {
    const { enabled } = this.args.config;
    return enabled === 'disable' || enabled === 'default-disabled';
  }

  get formattedStartDate() {
    const { startTimestamp } = this.args;
    return startTimestamp ? parseAPITimestamp(startTimestamp, 'MMMM yyyy') : null;
  }

  get formattedEndDate() {
    const { endTimestamp } = this.args;
    return endTimestamp ? parseAPITimestamp(endTimestamp, 'MMMM yyyy') : null;
  }

  // passed into page-header for the export modal alert
  get upgradesDuringActivity() {
    const { versionHistory, activity } = this.args;
    return filterVersionHistory(versionHistory, activity?.start_time, activity?.end_time);
  }

  get upgradeExplanations() {
    if (this.upgradesDuringActivity.length) {
      return this.upgradesDuringActivity.map((upgrade: VersionHistory) => {
        let explanation;
        const date = parseAPITimestamp(upgrade.timestamp_installed, 'MMM d, yyyy');
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
