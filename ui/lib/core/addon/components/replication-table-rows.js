import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-table-rows';

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
