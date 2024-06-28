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
  @tracked showCSVDownloadModal = false;

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

  destructureCountsToArray(object) {
    // destructure the namespace object  {label: 'some-namespace', entity_clients: 171, non_entity_clients: 20, acme_clients: 6, secret_syncs: 10, clients: 207}
    // to get integers for CSV file
    const { clients, entity_clients, non_entity_clients, acme_clients, secret_syncs } = object;
    const { isSecretsSyncActivated } = this.args;

    return [
      clients,
      entity_clients,
      non_entity_clients,
      acme_clients,
      ...(isSecretsSyncActivated ? [secret_syncs] : []),
    ];
  }

  constructCsvRow(namespaceColumn, mountColumn = null, totalColumns, newColumns = null) {
    // if namespaceColumn is a string, then we're at mount level attribution, otherwise it is an object
    // if constructing a namespace row, mountColumn=null so the column is blank, otherwise it is an object
    const otherColumns = newColumns ? [...totalColumns, ...newColumns] : [...totalColumns];
    return [
      `${typeof namespaceColumn === 'string' ? namespaceColumn : namespaceColumn.label}`,
      `${mountColumn ? mountColumn.label : '*'}`,
      ...otherColumns,
    ];
  }

  generateCsvData() {
    const totalAttribution = this.args.totalClientAttribution;
    const newAttribution = this.barChartNewClients ? this.args.newClientAttribution : null;
    const { isSecretsSyncActivated } = this.args;
    const csvData = [];
    // added to clarify that the row of namespace totals without an auth method (blank) are not additional clients
    // but indicate the total clients for that ns, including its auth methods
    const upgrade = this.args.upgradesDuringActivity?.length
      ? `\n **data contains an upgrade (mount summation may not equal namespace totals)`
      : '';
    const descriptionOfBlanks = this.isSingleNamespace
      ? ''
      : `\n *namespace totals, inclusive of mount clients${upgrade}`;
    // client type order here should match array order returned by destructureCountsToArray
    let csvHeader = [
      'Namespace path',
      `"Mount path${descriptionOfBlanks}"`, // double quotes necessary so description stays inside this cell
      'Total clients',
      'Entity clients',
      'Non-entity clients',
      'ACME clients',
      ...(isSecretsSyncActivated ? ['Secrets sync clients'] : []),
    ];

    if (newAttribution) {
      csvHeader = [
        ...csvHeader,
        'Total new clients',
        'New entity clients',
        'New non-entity clients',
        'New ACME clients',
        ...(isSecretsSyncActivated ? 'New secrets sync clients' : []),
      ];
    }

    totalAttribution.forEach((totalClientsObject) => {
      const namespace = this.isSingleNamespace ? this.args.selectedNamespace : totalClientsObject;
      const mount = this.isSingleNamespace ? totalClientsObject : null;

      // find new client data for namespace/mount object we're iterating over
      const newClientsObject = newAttribution
        ? newAttribution.find((d) => d.label === totalClientsObject.label)
        : null;

      const totalClients = this.destructureCountsToArray(totalClientsObject);
      const newClients = newClientsObject ? this.destructureCountsToArray(newClientsObject) : null;

      csvData.push(this.constructCsvRow(namespace, mount, totalClients, newClients));
      // constructCsvRow returns an array that corresponds to a row in the csv file:
      // ['ns label', 'mount label', total client #, entity #, non-entity #, acme #, secrets sync #, ...new client #'s]

      // only iterate through mounts if NOT viewing a single namespace
      if (!this.isSingleNamespace && namespace.mounts) {
        namespace.mounts.forEach((mount) => {
          const newMountData = newAttribution
            ? newClientsObject?.mounts.find((m) => m.label === mount.label)
            : null;
          const mountTotalClients = this.destructureCountsToArray(mount);
          const mountNewClients = newMountData ? this.destructureCountsToArray(newMountData) : null;
          csvData.push(this.constructCsvRow(namespace, mount, mountTotalClients, mountNewClients));
        });
      }
    });

    csvData.unshift(csvHeader);
    // make each nested array a comma separated string, join each array "row" in csvData with line break (\n)
    return csvData.map((d) => d.join()).join('\n');
  }

  get formattedCsvFileName() {
    const endRange = this.formattedEndDate ? `-${this.formattedEndDate}` : '';
    const csvDateRange = this.formattedStartDate ? `_${this.formattedStartDate + endRange}` : '';
    return this.isSingleNamespace
      ? `clients_by_mount_path${csvDateRange}`
      : `clients_by_namespace${csvDateRange}`;
  }

  get modalExportText() {
    const { isSecretsSyncActivated } = this.args;
    return `This export will include the namespace path, mount path and associated total entity, non-entity${
      isSecretsSyncActivated ? ', ACME and secrets sync clients' : ' and ACME clients'
    } for the ${this.formattedEndDate ? 'date range' : 'month'} below.`;
  }

  @action
  exportChartData(filename) {
    const contents = this.generateCsvData();
    this.download.csv(filename, contents);
    this.showCSVDownloadModal = false;
  }
}
