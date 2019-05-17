import Component from '@ember/component';
import { computed } from '@ember/object';

/**
 * @module VltTable
 * VltTable components render a table with the total number of HTTP Requests to a Vault server per month.
 *
 * @example
 * ```js
 * const DATA = [
 *  {
 *   foo: 'panda',
 *   bar: 50,
 *  }
 * ];
 *
 * <VltTable @data={{DATA}} />
 * ```
 *
 * @param data {Array} - A list of objects containing the table headers and corresponding values.
 */

export default Component.extend({
  classNames: ['vlt-table'],
  data: null,
});
