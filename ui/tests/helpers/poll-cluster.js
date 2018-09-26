import { run } from '@ember/runloop';
import { registerAsyncHelper } from '@ember/test';

export function pollCluster(owner) {
  const clusterRoute = owner.lookup('route:vault/cluster');
  return run(() => {
    return clusterRoute.controller.model.reload();
  });
}

registerAsyncHelper('pollCluster', function(app) {
  pollCluster(app.__container__);
});
