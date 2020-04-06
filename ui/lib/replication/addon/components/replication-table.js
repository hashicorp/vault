import Component from '@ember/component';
import layout from '../templates/components/replication-table';

/**
 * @module ReplicationTable
 * ReplicationTable components are used to...
 *
 * @example
 * ```js
 * <ReplicationTable @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

export default Component.extend({
  layout,
  data: null,
});
