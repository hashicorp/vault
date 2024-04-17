/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { next } from '@ember/runloop';
import { service } from '@ember/service';
import { computed } from '@ember/object';
import Controller, { inject as controller } from '@ember/controller';

export default Controller.extend({
  clusterController: controller('vault.cluster'),

  backendCrumb: computed('clusterController.model.name', function () {
    return {
      label: 'leases',
      text: 'leases',
      path: 'vault.cluster.access.leases.list-root',
      model: this.clusterController.model.name,
    };
  }),

  flashMessages: service(),
  router: service(),

  actions: {
    revokeLease(model) {
      return model.destroyRecord().then(() => {
        return this.router.transitionTo('vault.cluster.access.leases.list-root');
      });
    },

    renewLease(model, increment) {
      const adapter = model.store.adapterFor('lease');
      const flash = this.flashMessages;
      adapter
        .renew(model.id, increment?.seconds)
        .then(() => {
          this.send('refreshModel');
          // lol this is terrible, but there's no way to get the promise from the route refresh
          next(() => {
            flash.success(`The lease ${model.id} was successfully renewed.`);
          });
        })
        .catch((e) => {
          const errString = e.errors.join('.');
          flash.danger(`There was an error renewing the lease: ${errString}`);
        });
    },
  },
});
