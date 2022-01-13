import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

/**
 * @module Attribution
 * Attribution components display the top 10 total client counts for namespaces or auth methods (mounts) during a billing period.
 * If view is filtered for a single month, two graphs display and show a breakdown of new and total client counts by namespace or auth method, respectively 
 *
 * @example
 * ```js
 * <Clients::Attribution
 *    @newClientsData={{this.newClientData}}
 *    @totalClientsData={{this.totalClientData}}
 *    @chartLegend={{this.chartLegend}}
 *    @isDateRange={{this.isDateRange}}
 *    @isAllNamespaces={{this.isAllNamespaces}}
    >
      <button type="button"> Export attribution data </button>
    </Clients::Attribution>
 * ```

 * @param {array} newClientsData - (passed to child chart) must be an array of flattened objects
 * @param {array} totalClientsData - (passed to child chart) must be an array of flattened objects
 * @param {array} chartLegend - (passed to child) array of objects with key names 'key' and 'label' so data can be stacked
 * @param {boolean} isDateRange - discerns if dataset from API is a date range or single month
 * @param {boolean} isAllNamespaces - relays if filtered by all namespaces, or by single namespace
 * @param {string} activityDateRange - for modal and csv download to display correct date range
 */

export default class Attribution extends Component {
  @tracked showCSVDownloadModal = false;

  // TODO should this be a getter? It will be updated from parent component's filter bar
  get isDateRange() {
    return this.args.isDateRange;
  }

  get isAllNamespaces() {
    return this.args.isAllNamespaces;
  }

  get clientCountBreakdown() {
    return this.isAllNamespaces ? 'Namespace' : 'Auth method';
  }

  get chartText() {
    let dateText = this.isDateRange ? 'date range' : 'month';
    if (this.isAllNamespaces) {
      return {
        description:
          'This data shows the top ten namespaces by client count and can be used to understand where clients are originating. Namespaces are identified by path. To see all namespaces, export this data.',
        newCopy: `The new clients in the namespace for this ${dateText}. 
          This aids in understanding which namespaces create and use new clients 
          ${dateText === 'date range' ? ' over time.' : '.'}`,
        totalCopy: `The total clients in the namespace for this ${dateText}. This number is useful for identifying overall usage volume.`,
      };
    } else if (!this.isAllNamespaces) {
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

  get getCsvData() {
    let results = '',
      data,
      fields;
    if (this.isDateRange) {
      data = this.args.totalClientsData;
      fields = [`${this.clientCountBreakdown}`, 'Active clients', 'Unique entities', 'Non-entity tokens'];
    } else {
      data = this.args.newClientsData;
      // will need add a column for new clients
      fields = [`${this.clientCountBreakdown}`, 'Active clients', 'Unique entities', 'Non-entity tokens'];
    }

    results = fields.join(',') + '\n';
    data.forEach(function (item) {
      // debugger
      let path = item.label !== '' ? item.label : 'root',
        total = item.total,
        unique = item.distinct_entities,
        non_entity = item.non_entity_tokens;

      results += path + ',' + total + ',' + unique + ',' + non_entity + '\n';
    });
    return results;
  }
  // Return csv filename with start and end dates
  get getCsvFileName() {
    return this.args.activityDateRange
      ? `clients-by-${this.clientCountBreakdown}-${this.args.activityDateRange}`
      : `clients-by-${this.clientCountBreakdown}-${new Date()}`;
  }
}
