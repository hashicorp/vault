import Mixin from '@ember/object/mixin';
import utils from 'vault/lib/key-utils';
import { task } from 'ember-concurrency';

export default Mixin.create({
  navToNearestAncestor: task(function*(key) {
    let ancestors = utils.ancestorKeysForKey(key);
    let errored = false;
    let nearest = ancestors.pop();
    while (nearest) {
      try {
        let transition = this.transitionToRoute('vault.cluster.secrets.backend.list', nearest);
        transition.data.isDeletion = true;
        yield transition.promise;
      } catch (e) {
        errored = true;
        nearest = ancestors.pop();
      } finally {
        if (!errored) {
          nearest = null;
          // eslint-disable-next-line
          return;
        }
        errored = false;
      }
    }
    this.transitionToRoute('vault.cluster.secrets.backend.list-root');
  }),
});
