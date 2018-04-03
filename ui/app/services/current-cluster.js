import Ember from 'ember';

export default Ember.Service.extend({
  cluster: null,

  setCluster(cluster) {
    this.set('cluster', cluster);
  },
});
