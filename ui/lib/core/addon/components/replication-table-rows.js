import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-table-rows';

export default Component.extend({
  layout,
  classNames: ['replication-table-rows'],
  data: null,
  clusterDetails: computed('data', function() {
    const { data } = this;
    return data.dr || data;
  }),
  mode: computed('clusterDetails', function() {
    return this.clusterDetails.mode || 'unknown';
  }),
  merkleRoot: computed('clusterDetails', function() {
    return this.clusterDetails.merkleRoot || 'unknown';
  }),
  clusterId: computed('clusterDetails', function() {
    return this.clusterDetails.clusterId || 'unknown';
  }),
  syncProgress: computed('clusterDetails', function() {
    return this.clusterDetails.syncProgress || false;
  }),
});
