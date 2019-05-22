import ClusterRouteBase from './cluster-route-base';

export default ClusterRouteBase.extend({
  model() {
    return this.store.queryRecord('requests', {});
  },
});
