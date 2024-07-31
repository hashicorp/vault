/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { sanitizePath } from 'core/utils/sanitize-path';
import { format, isSameMonth } from 'date-fns';
import { task } from 'ember-concurrency';

/**
 * @module ClientsExportButtonComponent
 * ClientsExportButton components are used to display the export button, manage the modal, and download the file from the clients export API
 *
 * @example
 * ```js
 * <Clients::ExportButton @startTimestamp="2022-06-01T23:00:11.050Z" @endTimestamp="2022-12-01T23:00:11.050Z" @selectedNamespace="foo" />
 * ```
 * @param {string} [startTimestamp] - ISO timestamp of start time, to be passed to export request
 * @param {string} [endTimestamp] - ISO timestamp of end time, to be passed to export request
 * @param {string} [namespace] - namespace filter. Will be appended to the current namespace in the export request.
 */
export default class ClientsExportButtonComponent extends Component {
  @service download;
  @service namespace;
  @service store;

  @tracked showExportModal = false;
  @tracked exportFormat = 'csv';
  @tracked downloadError = '';

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
    const { selectedNamespace } = this.args;
    return selectedNamespace ? sanitizePath(`${currentNs}/${selectedNamespace}`) : sanitizePath(currentNs);
  }

  async getExportData() {
    const adapter = this.store.adapterFor('clients/activity');
    const { startTimestamp, endTimestamp } = this.args;
    return adapter.exportData({
      format: this.exportFormat,
      start_time: startTimestamp,
      end_time: endTimestamp,
      namespace: this.namespaceFilter,
    });
  }

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
