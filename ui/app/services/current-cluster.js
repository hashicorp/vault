import Service from '@ember/service';

export default Service.extend({
  cluster: null,

  setCluster(cluster) {
    this.set('cluster', cluster);
  },
});
