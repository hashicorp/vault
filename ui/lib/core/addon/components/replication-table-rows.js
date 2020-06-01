import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-table-rows';

/**
 * @module ReplicationTableRows
 * The `ReplicationTableRows` component is table component.  It displays cluster mode details specific to the cluster of the Dashboard it is used on.
 *
 * @example
 * ```js
 * <ReplicationTableRows
    @replicationDetails={{replicationDetails}}
    @clusterMode="primary"
    />
 * ```
 * @param {Object} replicationDetails=null - An Ember data object pulled from the Ember Model. It contains details specific to the whether the replication is dr or performance.
 * @param {String} clusterMode=null - The cluster mode passed through to a table component. 
 */

export default Component.extend({
  layout,
  classNames: ['replication-table-rows'],
  replicationDetails: null,
  clusterMode: null,
  merkleRoot: computed('replicationDetails.{merkleRoot}', function() {
    return this.replicationDetails.merkleRoot || 'unknown';
  }),
  clusterId: computed('replicationDetails.{clusterId}', function() {
    return this.replicationDetails.clusterId || 'unknown';
  }),
});
