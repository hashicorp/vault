/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { format, isSameMonth } from 'date-fns';

/**
 * @module Attribution
 * Attribution components display the top 10 total client counts for namespaces or auth methods (mounts) during a billing period.
 * A horizontal bar chart shows on the right, with the top namespace/auth method and respective client totals on the left.
 *
 * @example
 *  <Clients::Attribution
 *    @totalUsageCounts={{this.totalUsageCounts}}
 *    @newUsageCounts={{this.newUsageCounts}}
 *    @totalClientAttribution={{this.totalClientAttribution}}
 *    @newClientAttribution={{this.newClientAttribution}}
 *    @selectedNamespace={{this.selectedNamespace}}
 *    @startTimestamp={{this.startTime}}
 *    @endTimestamp={{this.endTime}}
 *    @isHistoricalMonth={{false}}
 *    @responseTimestamp={{this.responseTimestamp}}
 *    @upgradesDuringActivity={{array (hash version="1.10.1" previousVersion="1.9.1" timestampInstalled= "2021-11-18T10:23:16Z") }}
 *  />
 *
 * @param {object} totalUsageCounts - object with total client counts for chart tooltip text
 * @param {object} newUsageCounts - object with new client counts for chart tooltip text
 * @param {array} totalClientAttribution - array of objects containing a label and breakdown of client counts for total clients
 * @param {array} newClientAttribution - array of objects containing a label and breakdown of client counts for new clients
 * @param {string} selectedNamespace - namespace selected from filter bar
 * @param {string} startTimestamp - timestamp string from activity response to render start date for CSV modal and whether copy reads 'month' or 'date range'
 * @param {string} endTimestamp - timestamp string from activity response to render end date for CSV modal and whether copy reads 'month' or 'date range'
 * @param {string} responseTimestamp -  ISO timestamp created in serializer to timestamp the response, renders in bottom left corner below attribution chart
 * @param {boolean} isHistoricalMonth - when true data is from a single, historical month so side-by-side charts should display for attribution data
 * @param {array} upgradesDuringActivity - array of objects containing version history upgrade data
 * @param {boolean} isSecretsSyncActivated - boolean to determine if secrets sync is activated
 */

export default class Attribution extends Component {
  @service download;
  @service store;
  @tracked showCSVDownloadModal = false;
  @tracked downloadError = '';

  get attributionLegend() {
    const attributionLegend = [
      { key: 'entity_clients', label: 'entity clients' },
      { key: 'non_entity_clients', label: 'non-entity clients' },
      { key: 'acme_clients', label: 'ACME clients' },
    ];

    if (this.args.isSecretsSyncActivated) {
      attributionLegend.push({ key: 'secret_syncs', label: 'secrets sync clients' });
    }
    return attributionLegend;
  }

  get formattedStartDate() {
    if (!this.args.startTimestamp) return null;
    return parseAPITimestamp(this.args.startTimestamp, 'MMMM yyyy');
  }

  get formattedEndDate() {
    if (!this.args.startTimestamp && !this.args.endTimestamp) return null;
    // displays on CSV export modal, no need to display duplicate months and years
    const startDateObject = parseAPITimestamp(this.args.startTimestamp);
    const endDateObject = parseAPITimestamp(this.args.endTimestamp);
    return isSameMonth(startDateObject, endDateObject) ? null : format(endDateObject, 'MMMM yyyy');
  }

  get hasCsvData() {
    return this.args.totalClientAttribution ? this.args.totalClientAttribution.length > 0 : false;
  }

  get isSingleNamespace() {
    // if a namespace is selected, then we're viewing top 10 auth methods (mounts)
    return !!this.args.selectedNamespace;
  }

  // truncate data before sending to chart component
  get barChartTotalClients() {
    return this.args.totalClientAttribution?.slice(0, 10);
  }

  get barChartNewClients() {
    return this.args.newClientAttribution?.slice(0, 10);
  }

  get topClientCounts() {
    // get top namespace or auth method
    return this.args.totalClientAttribution ? this.args.totalClientAttribution[0] : null;
  }

  get attributionBreakdown() {
    // display text for hbs
    return this.isSingleNamespace ? 'auth method' : 'namespace';
  }

  get chartText() {
    if (!this.args.totalClientAttribution) {
      return { description: 'There is a problem gathering data' };
    }
    const dateText = this.formattedEndDate ? 'date range' : 'month';
    switch (this.isSingleNamespace) {
      case true:
        return {
          description:
            'This data shows the top ten authentication methods by client count within this namespace, and can be used to understand where clients are originating. Authentication methods are organized by path.',
          newCopy: `The new clients used by the auth method for this ${dateText}. This aids in understanding which auth methods create and use new clients${
            dateText === 'date range' ? ' over time.' : '.'
          }`,
          totalCopy: `The total clients used by the auth method for this ${dateText}. This number is useful for identifying overall usage volume. `,
        };
      case false:
        return {
          description:
            'This data shows the top ten namespaces by client count and can be used to understand where clients are originating. Namespaces are identified by path. To see all namespaces, export this data.',
          newCopy: `The new clients in the namespace for this ${dateText}.
          This aids in understanding which namespaces create and use new clients${
            dateText === 'date range' ? ' over time.' : '.'
          }`,
          totalCopy: `The total clients in the namespace for this ${dateText}. This number is useful for identifying overall usage volume.`,
        };
      default:
        return '';
    }
  }

  async generateCsvData() {
    const adapter = this.store.adapterFor('clients/activity');
    const { startTimestamp, endTimestamp } = this.args;
    return adapter.exportData({
      start_time: startTimestamp,
      end_time: endTimestamp,
    });
  }

  get formattedCsvFileName() {
    const endRange = this.formattedEndDate ? `-${this.formattedEndDate}` : '';
    const csvDateRange = this.formattedStartDate ? `_${this.formattedStartDate + endRange}` : '';
    return this.isSingleNamespace
      ? `clients_by_mount_path${csvDateRange}`
      : `clients_by_namespace${csvDateRange}`;
  }

  @action
  async exportChartData(filename) {
    try {
      const contents = await this.generateCsvData();
      this.download.csv(filename, contents);
      this.showCSVDownloadModal = false;
    } catch (e) {
      this.downloadError = e.message;
    }
  }

  @action resetModal() {
    this.showCSVDownloadModal = false;
    this.downloadError = '';
  }
}
