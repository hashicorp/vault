import Ember from 'ember';

export default Ember.Controller.extend({
  auth: Ember.inject.service(),

  actions: {
    seal() {
      return this.model.cluster.store.adapterFor('cluster').seal().then(() => {
        this.model.cluster.get('leaderNode').set('sealed', true);
        this.get('auth').deleteCurrentToken();
        return this.transitionToRoute('vault.cluster');
      });
    },
  },
});
