/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { next } from '@ember/runloop';
import { service } from '@ember/service';
import { computed } from '@ember/object';
import Controller, { inject as controller } from '@ember/controller';

export default Controller.extend({
  clusterController: controller('vault.cluster'),
  api: service(),
  flashMessages: service(),
  router: service(),

  backendCrumb: computed('clusterController.model.name', function () {
    return [
      {
        label: 'Vault',
        text: 'Vault',
        path: 'vault.cluster.dashboard',
        icon: 'vault',
      },
      {
        label: 'Leases',
        text: 'Leases',
        path: 'vault.cluster.access.leases.list-root',
        model: this.clusterController.model.name,
      },
    ];
  }),

  actions: {
    revokeLease(model) {
      return this.api.sys
        .leasesRevokeLease({ lease_id: model.id })
        .then(() => {
          return this.router.transitionTo('vault.cluster.access.leases.list-root');
        })
        .catch((e) => {
          this.api.parseError(e).then(({ message }) => {
            this.flashMessages.danger(`There was an error revoking the lease: ${message}`);
          });
        });
    },

    renewLease(model, increment) {
      const flash = this.flashMessages;
      this.api.sys
        .leasesRenewLease({ lease_id: model.id, increment: String(increment?.seconds ?? '') })
        .then(() => {
          this.send('refreshModel');
          // lol this is terrible, but there's no way to get the promise from the route refresh
          next(() => {
            flash.success(`The lease ${model.id} was successfully renewed.`);
          });
        })
        .catch((e) => {
          this.api.parseError(e).then(({ message }) => {
            flash.danger(`There was an error renewing the lease: ${message}`);
          });
        });
    },
  },
});
