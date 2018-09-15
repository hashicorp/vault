import { run } from '@ember/runloop';
import { registerAsyncHelper } from '@ember/test';

export default registerAsyncHelper('pollCluster', function(app) {
  const clusterRoute = app.__container__.cache['route:vault/cluster'];
  return run(() => {
    return clusterRoute.controller.model.reload();
  });
});
