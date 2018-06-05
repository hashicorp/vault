import Ember from 'ember';

export default Ember.Test.registerAsyncHelper('pollCluster', function(app) {
  const clusterRoute = app.__container__.cache['route:vault/cluster'];
  return Ember.run(() => {
    return clusterRoute.controller.model.reload();
  });
});
