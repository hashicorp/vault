import { next } from '@ember/runloop';
import { inject as service } from '@ember/service';
import { computed } from '@ember/object';
import Controller, { inject as controller } from '@ember/controller';

export default Controller.extend({
  clusterController: controller('vault.cluster'),

  backendCrumb: computed(function() {
    return {
      label: 'leases',
      text: 'leases',
      path: 'vault.cluster.access.leases.list-root',
      model: this.get('clusterController.model.name'),
    };
  }),

  flashMessages: service(),

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
          next(() => {
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
