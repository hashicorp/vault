import { set } from '@ember/object';
import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import UnloadModelRoute from 'vault/mixins/unload-model-route';

import utils from 'vault/lib/key-utils';

export default Route.extend(UnloadModelRoute, {
  beforeModel() {
    const { lease_id: leaseId } = this.paramsFor(this.routeName);
    const parentKey = utils.parentKeyForKey(leaseId);
    if (utils.keyIsFolder(leaseId)) {
      if (parentKey) {
        return this.transitionTo('vault.cluster.access.leases.list', parentKey);
      } else {
        return this.transitionTo('vault.cluster.access.leases.list-root');
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
