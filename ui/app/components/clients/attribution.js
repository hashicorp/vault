import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

/**
 * @module Attribution
 * Attribution components display the top 10 total client counts for namespaces or auth methods (mounts) during a billing period.
 * A horizontal bar chart shows on the right, with the top namespace/auth method and respective client totals on the left.
 *
 * @example
 * ```js
 *  <Clients::Attribution
 *    @chartLegend={{this.chartLegend}}
 *    @topTenNamespaces={{this.topTenNamespaces}}
 *    @runningTotals={{this.runningTotals}}
 *    @selectedNamespace={{this.selectedNamespace}}
 *    @startTimeDisplay={{this.startTimeDisplay}}
 *    @endTimeDisplay={{this.endTimeDisplay}}
 *    @isDateRange={{this.isDateRange}}
 *    @timestamp={{this.responseTimestamp}}
 *  />
 * ```
 * @param {array} chartLegend - (passed to child) array of objects with key names 'key' and 'label' so data can be stacked
 * @param {array} topTenNamespaces - (passed to child chart) array of top 10 namespace objects
 * @param {object} runningTotals - object with total client counts for chart tooltip text
 * @param {string} selectedNamespace - namespace selected from filter bar
 * @param {string} startTimeDisplay - start date for CSV modal
 * @param {string} endTimeDisplay - end date for CSV modal
 * @param {boolean} isDateRange - getter calculated in parent to relay if dataset is a date range or single month
 * @param {string} timestamp - timestamp response was received from API
 */

export default class Attribution extends Component {
  @tracked showCSVDownloadModal = false;

  get isDateRange() {
    return this.args.isDateRange;
  }

  get isSingleNamespace() {
    // if a namespace is selected, then we're viewing top 10 auth methods (mounts)
    return !!this.args.selectedNamespace;
  }

  get totalClientsData() {
    // get dataset for bar chart displaying top 10 namespaces/mounts with highest # of total clients
    return this.isSingleNamespace
      ? this.filterByNamespace(this.args.selectedNamespace)
      : this.args.topTenNamespaces;
  }

  get topClientCounts() {
    // get top namespace or auth method
    return this.totalClientsData[0];
  }

  get attributionBreakdown() {
    // display 'Auth method' or 'Namespace' respectively in CSV file
    return this.isSingleNamespace ? 'Auth method' : 'Namespace';
  }

  get chartText() {
    let dateText = this.isDateRange ? 'date range' : 'month';
    if (!this.isSingleNamespace) {
      return {
        description:
          'This data shows the top ten namespaces by client count and can be used to understand where clients are originating. Namespaces are identified by path. To see all namespaces, export this data.',
        newCopy: `The new clients in the namespace for this ${dateText}. 
          This aids in understanding which namespaces create and use new clients 
          ${dateText === 'date range' ? ' over time.' : '.'}`,
        totalCopy: `The total clients in the namespace for this ${dateText}. This number is useful for identifying overall usage volume.`,
      };
    } else if (this.isSingleNamespace) {
      return {
        description:
          'This data shows the top ten authentication methods by client count within this namespace, and can be used to understand where clients are originating. Authentication methods are organized by path.',
        newCopy: `The new clients used by the auth method for this ${dateText}. This aids in understanding which auth methods create and use new clients 
        ${dateText === 'date range' ? ' over time.' : '.'}`,
        totalCopy: `The total clients used by the auth method for this ${dateText}. This number is useful for identifying overall usage volume. `,
      };
    } else {
      return {
        description: 'There is a problem gathering data',
        newCopy: 'There is a problem gathering data',
        totalCopy: 'There is a problem gathering data',
      };
    }
  }

  // TODO CMB update with proper data format when we have
  get getCsvData() {
    let results = '',
      data,
      fields;

    // TODO CMB will CSV for namespaces include mounts?
    fields = [`${this.attributionBreakdown}`, 'Active clients', 'Unique entities', 'Non-entity tokens'];

    results = fields.join(',') + '\n';
    data.forEach(function (item) {
      let path = item.label !== '' ? item.label : 'root',
        total = item.total,
        unique = item.entity_clients,
        non_entity = item.non_entity_clients;

      results += path + ',' + total + ',' + unique + ',' + non_entity + '\n';
    });
    return results;
  }
  // TODO CMB - confirm with design file name structure
  get getCsvFileName() {
    let activityDateRange = `${this.args.startTimeDisplay} - ${this.args.endTimeDisplay}`;
    return activityDateRange
      ? `clients-by-${this.attributionBreakdown}-${activityDateRange}`
      : `clients-by-${this.attributionBreakdown}-${new Date()}`;
  }

  // HELPERS
  filterByNamespace(namespace) {
    // return top 10 mounts for a namespace
    let namespaceObject = this.args.topTenNamespaces.find((ns) => ns.label === `${namespace}/`);
    // debugger;
    if (namespaceObject.mounts) {
      this.noMountsOnPayload = true;
      return namespaceObject.mounts.slice(0, 10);
    } else {
      this.noMountsOnPayload = false;
      // ARG TODO unsure on what to return here?
    }
  }
}
