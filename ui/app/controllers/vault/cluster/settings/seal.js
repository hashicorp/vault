import { inject as service } from '@ember/service';
import Controller from '@ember/controller';

export default Controller.extend({
  auth: service(),

  actions: {
    seal() {
      return this.model.cluster.store
        .adapterFor('cluster')
        .seal()
        .then(() => {
          this.model.cluster.get('leaderNode').set('sealed', true);
          this.get('auth').deleteCurrentToken();
          return this.transitionToRoute('vault.cluster');
        });
    },
  },
});
