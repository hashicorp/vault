/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import { waitFor } from '@ember/test-waiters';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { sanitizePath } from 'core/utils/sanitize-path';
import { task } from 'ember-concurrency';
import { formatExportData } from 'core/utils/client-counts/serializers';

import type DownloadService from 'vault/services/download';
import type FlagsService from 'vault/services/flags';
import type NamespaceService from 'vault/services/namespace';
import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type VersionService from 'vault/services/version';
import type CapabilitiesService from 'vault/services/capabilities';
import type Owner from '@ember/owner';
import type { HTMLElementEvent } from 'vault/forms';
import type { Extensions } from 'vault/services/download';

/**
 * @module ClientsPageHeader
 * ClientsPageHeader components are used to render a header and check for export capabilities before rendering an export button.
 *
 * @example
 * ```js
 * <Clients::PageHeader @startTimestamp="2022-06-01T23:00:11.050Z" @endTimestamp="2022-12-01T23:00:11.050Z" @namespace="foo" @upgradesDuringActivity={{array (hash version="1.10.1" previousVersion="1.9.1" timestampInstalled= "2021-11-18T10:23:16Z") }} />
 * ```
 * @param {Date} [billingStartTime] - billing start date, to be passed to date picker
 * @param {string} [activityTimestamp] -  ISO timestamp created in serializer to timestamp the response to be displayed in page header
 * @param {Date} [startTimestamp] - start time, to be passed to export request
 * @param {Date} [endTimestamp] - end time, to be passed to export request
 * @param {number} [retentionMonths = 48] - number of months for historical billing, to be passed to date picker
 * @param {string} [upgradesDuringActivity] - array of objects containing version history upgrade data
 * @param {boolean} [noData = false] - when true, export button will hide regardless of capabilities
 * @param {function} [onChange] - callback when a new date range is saved, to be passed to date picker
 */

interface Args {
  billingStartTime: Date;
  retentionMonths: number;
  activityTimestamp: string;
  startTimestamp: Date;
  endTimestamp: Date;
  upgradesDuringActivity: string[];
  noData: boolean;
  onChange: CallableFunction;
}

export default class ClientsPageHeaderComponent extends Component<Args> {
  @service declare readonly download: DownloadService;
  @service declare readonly flags: FlagsService;
  @service declare readonly namespace: NamespaceService;
  @service declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly version: VersionService;
  @service declare readonly capabilities: CapabilitiesService;

  @tracked canDownload = false;
  @tracked showEditModal = false;
  @tracked showExportModal = false;
  @tracked exportFormat: keyof Extensions = 'csv';
  @tracked downloadError = '';

  constructor(owner: Owner, args: Args) {
    super(owner, args);
    this.getExportCapabilities();
  }

  get showExportButton() {
    if (this.args.noData === true) return false;
    return this.canDownload;
  }

  @waitFor
  async getExportCapabilities() {
    const ns = this.namespace.path;
    try {
      // selected namespace usually ends in /
      const namespace = sanitizePath(ns);
      const { canSudo } = await this.capabilities.for('clientsActivityExport', { namespace });
      this.canDownload = canSudo;
    } catch (e) {
      // if we can't read capabilities, default to show
      this.canDownload = true;
    }
  }

  get formattedStartDate() {
    return this.args.startTimestamp ? parseAPITimestamp(this.args.startTimestamp, 'MMMM yyyy') : null;
  }

  get formattedEndDate() {
    return this.args.endTimestamp ? parseAPITimestamp(this.args.endTimestamp, 'MMMM yyyy') : null;
  }

  get showEndDate() {
    // displays on CSV export modal, no need to display duplicate months and years
    return this.formattedEndDate && this.formattedStartDate !== this.formattedEndDate;
  }

  get formattedCsvFileName() {
    const endRange = this.showEndDate ? `-${this.formattedEndDate}` : '';
    const csvDateRange = this.formattedStartDate ? `_${this.formattedStartDate + endRange}` : '';
    const ns = this.namespace.path ? `_${this.namespace.path}` : '';
    return `clients_export${ns}${csvDateRange}`;
  }

  async getExportData() {
    const { startTimestamp, endTimestamp } = this.args;
    const namespace = this.namespace.path;
    const headers = namespace ? this.api.buildHeaders({ namespace }) : undefined;
    const { raw } = await this.api.sys.internalClientActivityExportRaw(
      {
        // the API only accepts json or csv
        format: this.exportFormat === 'jsonl' ? 'json' : 'csv',
        start_time: startTimestamp ? startTimestamp.toISOString() : undefined,
        end_time: endTimestamp ? endTimestamp.toISOString() : undefined,
      },
      headers
    );
    if (raw.status !== 200) {
      throw { message: 'No data to export in provided time range.' };
    }
    return formatExportData(raw, { isDownload: true });
  }

  exportChartData = task({ drop: true }, async (filename) => {
    try {
      const contents = await this.getExportData();
      this.download.download(filename, contents, this.exportFormat);
      this.showExportModal = false;
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.downloadError = message;
    }
  });

  @action
  refreshRoute() {
    this.router.refresh(this.router.currentRoute?.parent?.name);
  }

  @action
  resetModal() {
    this.showExportModal = false;
    this.downloadError = '';
  }

  @action
  setEditModalVisible(visible: boolean) {
    this.showEditModal = visible;
  }

  @action
  setExportFormat(event: HTMLElementEvent<HTMLInputElement>) {
    const { value } = event.target;
    this.exportFormat = value as keyof Extensions;
  }

  // LOCAL TEMPLATE HELPERS
  parseAPITimestamp = (timestamp: string, format: string) => {
    return parseAPITimestamp(timestamp, format);
  };
}
