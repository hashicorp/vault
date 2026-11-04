/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import UnloadModel from 'vault/mixins/unload-model-route';

export default Route.extend(UnloadModel, {
  capabilities: service(),
  store: service(),
  version: service(),

  beforeModel() {
    return this.version.fetchFeatures().then(() => {
      return this._super(...arguments);
    });
  },

  model(params) {
    return this.version.hasControlGroups ? this.store.findRecord('control-group', params.accessor) : null;
  },

  // The requesting entity and approvers link to the entity details page only when the user can read them
  async afterModel(model) {
    const ids = new Set((model?.authorizations || []).map((authz) => authz.id));
    if (model?.requestEntity?.id) ids.add(model.requestEntity.id);

    const pathById = {};
    for (const id of ids) {
      pathById[id] = this.capabilities.pathFor('identityCapabilities', { identityType: 'entity', id });
    }
    const paths = Object.values(pathById);
    const capabilities = paths.length ? await this.capabilities.fetch(paths) : {};
    this.canReadEntity = {};
    for (const id of ids) {
      this.canReadEntity[id] = Boolean(capabilities[pathById[id]]?.canRead);
    }
  },

  setupController(controller) {
    this._super(...arguments);
    controller.set('canReadEntity', this.canReadEntity);
  },

  actions: {
    willTransition() {
      return true;
    },
    // deactivate happens later than willTransition,
    // so since we're using the model to render links
    // we don't want the UI blinking
    deactivate() {
      this.unloadModel();
      return true;
    },
  },
});
