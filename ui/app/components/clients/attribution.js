import Component from '@glimmer/component';

// TODO: fill out below!!
/**
 * @module Attribution
 * Attribution components are used to...
 *
 * @example
 * ```js
 * <Attribution @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}>
 * Pass in export button
 * </Attribution>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

export default class Attribution extends Component {
  get dateRange() {
    // some conditional that returns "date range" or "month" depending on what the params are
    return 'date range';
  }

  get chartText() {
    // something that determines if data is by namespace or by auth method
    // and returns text
    // if byNamespace
    return {
      description:
        'This data shows the top ten namespaces by client count and can be used to understand where clients are originating. Namespaces are identified by path. To see all namespaces, export this data.',
      newCopy: `The new clients in the namespace for this ${this.dateRange}. 
        This aids in understanding which namespaces create and use new clients 
        ${this.dateRange === 'date range' ? ' over time.' : '.'}`,
      totalCopy: `The total clients in the namespace for this ${this.dateRange}. This number is useful for identifying overall usage volume.`,
    };
    // if byAuthMethod
    // return
    // byAuthMethod = {
    //   description: "This data shows the top ten authentication methods by client count within this namespace, and can be used to understand where new clients and total clients are originating. Authentication methods are organized by path.",
    //   newCopy: `The new clients used by the auth method for this {{@range}}. This aids in understanding which auth methods create and use new clients ${this.dateRange === "date range" ? " over time." : "."}`,
    //   totalCopy: `The total clients used by the auth method for this ${this.dateRange}. This number is useful for identifying overall usage volume. `
    // }
  }
}
