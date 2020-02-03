import { run } from '@ember/runloop';

export function pollCluster(owner) {
  const store = owner.lookup('service:store');
  return run(() => {
    return store.peekAll('cluster').firstObject.reload();
  });
}
