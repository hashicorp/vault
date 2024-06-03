/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { hash } from 'rsvp';
import { set } from '@ember/object';
import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { keyIsFolder, parentKeyForKey } from 'core/utils/key-utils';
import UnloadModelRoute from 'vault/mixins/unload-model-route';

export default Route.extend(UnloadModelRoute, {
  store: service(),
  router: service(),

  beforeModel() {
    const { lease_id: leaseId } = this.paramsFor(this.routeName);
    const parentKey = parentKeyForKey(leaseId);
    if (keyIsFolder(leaseId)) {
      if (parentKey) {
        return this.router.transitionTo('vault.cluster.access.leases.list', parentKey);
      } else {
        return this.router.transitionTo('vault.cluster.access.leases.list-root');
      }
    }
  },

  model(params) {
    const { lease_id } = params;
    return hash({
      lease: this.store.queryRecord('lease', {
        lease_id,
      }),
      capabilities: hash({
        renew: this.store.findRecord('capabilities', 'sys/leases/renew'),
        revoke: this.store.findRecord('capabilities', 'sys/leases/revoke'),
        leases: this.modelFor('vault.cluster.access.leases'),
      }),
    });
  },

  setupController(controller, model) {
    this._super(...arguments);
    const { lease_id: leaseId } = this.paramsFor(this.routeName);
    controller.setProperties({
      model: model.lease,
      capabilities: model.capabilities,
      baseKey: { id: leaseId },
    });
  },

  actions: {
    error(error) {
      const { lease_id } = this.paramsFor(this.routeName);
      set(error, 'keyId', lease_id);
      return true;
    },

    refreshModel() {
      this.refresh();
    },
  },
});
