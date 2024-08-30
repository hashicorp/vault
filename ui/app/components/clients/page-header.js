/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import { waitFor } from '@ember/test-waiters';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { sanitizePath } from 'core/utils/sanitize-path';
import { format, isSameMonth } from 'date-fns';
import { task } from 'ember-concurrency';

/**
 * @module ClientsPageHeader
 * ClientsPageHeader components are used to render a header and check for export capabilities before rendering an export button.
 *
 * @example
 * ```js
 * <Clients::PageHeader @startTimestamp="2022-06-01T23:00:11.050Z" @endTimestamp="2022-12-01T23:00:11.050Z" @namespace="foo" @upgradesDuringActivity={{array (hash version="1.10.1" previousVersion="1.9.1" timestampInstalled= "2021-11-18T10:23:16Z") }} />
 * ```
 * @param {string} [startTimestamp] - ISO timestamp of start time, to be passed to export request
 * @param {string} [endTimestamp] - ISO timestamp of end time, to be passed to export request
 * @param {string} [namespace] - namespace filter. Will be appended to the current namespace in the export request.
 * @param {string} [upgradesDuringActivity] - array of objects containing version history upgrade data
 * @param {boolean} [noData = false] - when true, export button will hide regardless of capabilities
 */
export default class ClientsPageHeaderComponent extends Component {
  @service download;
  @service namespace;
  @service store;

  @tracked canDownload = false;
  @tracked showExportModal = false;
  @tracked exportFormat = 'csv';
  @tracked downloadError = '';

  constructor() {
    super(...arguments);
    this.getExportCapabilities(this.args.namespace);
  }

  get showExportButton() {
    if (this.args.noData === true) return false;
    return this.canDownload;
  }

  @waitFor
  async getExportCapabilities(ns = '') {
    try {
      // selected namespace usually ends in /
      const url = ns
        ? `${sanitizePath(ns)}/sys/internal/counters/activity/export`
        : 'sys/internal/counters/activity/export';
      const cap = await this.store.findRecord('capabilities', url);
      this.canDownload = cap.canSudo;
    } catch (e) {
      // if we can't read capabilities, default to show
      this.canDownload = true;
    }
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

  get formattedCsvFileName() {
    const endRange = this.formattedEndDate ? `-${this.formattedEndDate}` : '';
    const csvDateRange = this.formattedStartDate ? `_${this.formattedStartDate + endRange}` : '';
    const ns = this.namespaceFilter ? `_${this.namespaceFilter}` : '';
    return `clients_export${ns}${csvDateRange}`;
  }

  get namespaceFilter() {
    const currentNs = this.namespace.path;
    const { namespace } = this.args;
    return namespace ? sanitizePath(`${currentNs}/${namespace}`) : sanitizePath(currentNs);
  }

  async getExportData() {
    const adapter = this.store.adapterFor('clients/activity');
    const { startTimestamp, endTimestamp } = this.args;
    return adapter.exportData({
      // the API only accepts json or csv
      format: this.exportFormat === 'jsonl' ? 'json' : 'csv',
      start_time: startTimestamp,
      end_time: endTimestamp,
      namespace: this.namespaceFilter,
    });
  }

  parseAPITimestamp = (timestamp, format) => {
    return parseAPITimestamp(timestamp, format);
  };

  exportChartData = task({ drop: true }, async (filename) => {
    try {
      const contents = await this.getExportData();
      this.download.download(filename, contents, this.exportFormat);
      this.showExportModal = false;
    } catch (e) {
      this.downloadError = e.message;
    }
  });

  @action setExportFormat(evt) {
    const { value } = evt.target;
    this.exportFormat = value;
  }

  @action resetModal() {
    this.showExportModal = false;
    this.downloadError = '';
  }
}
