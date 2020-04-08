import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-table-rows';

export default Component.extend({
  layout,
  data: null,
  clusterDetails: computed('data', function() {
    const { data } = this.data;
    return data.dr || data;
  }),
  mode: computed('clusterDetails', function() {
    const { clusterDetails } = this;
    return clusterDetails.mode || 'unknown';
  }),
  merkleRoot: computed('clusterDetails', function() {
    const { clusterDetails } = this;
    return clusterDetails.merkleRoot || 'unknown';
  }),
  clusterId: computed('clusterDetails', function() {
    const { clusterDetails } = this;
    return clusterDetails.clusterId || 'unknown';
  }),
  syncProgress: computed('clusterDetails', function() {
    const { clusterDetails } = this;
    return clusterDetails.syncProgress || false;
  }),
});
