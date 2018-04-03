import Ember from 'ember';

export default Ember.Controller.extend({
  clusterController: Ember.inject.controller('vault.cluster'),

  backendCrumb: Ember.computed(function() {
    return {
      label: 'leases',
      text: 'leases',
      path: 'vault.cluster.access.leases.list-root',
      model: this.get('clusterController.model.name'),
    };
  }),

  flashMessages: Ember.inject.service(),

  actions: {
    revokeLease(model) {
      return model.destroyRecord().then(() => {
        return this.transitionToRoute('vault.cluster.access.leases.list-root');
      });
    },

    renewLease(model, interval) {
      const adapter = model.store.adapterFor('lease');
      const flash = this.get('flashMessages');
      adapter
        .renew(model.id, interval)
        .then(() => {
          this.send('refreshModel');
          // lol this is terrible, but there's no way to get the promise from the route refresh
          Ember.run.next(() => {
            flash.success(`The lease ${model.id} was successfully renewed.`);
          });
        })
        .catch(e => {
          const errString = e.errors.join('.');
          flash.danger(`There was an error renewing the lease: ${errString}`);
        });
    },
  },
});
