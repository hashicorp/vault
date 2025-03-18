/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import UnloadModel from 'vault/mixins/unload-model-route';
import { service } from '@ember/service';
import { ROUTES } from 'vault/utils/routes';

export default Route.extend(UnloadModel, {
  router: service(),
  store: service(),
  templateName: 'vault/cluster/secrets/backend/sign',

  backendModel() {
    return this.modelFor('vault.cluster.secrets.backend');
  },

  pathQuery(role, backend) {
    return {
      id: `${backend}/sign/${role}`,
    };
  },

  pathForType() {
    return 'sign';
  },

  model(params) {
    const role = params.secret;
    const backendModel = this.backendModel();
    const backend = backendModel.id;

    if (backendModel.type !== 'ssh') {
      return this.router.transitionTo(ROUTES.VAULT_CLUSTER_SECRETS_BACKEND_LISTROOT, backend);
    }
    return this.store.queryRecord('capabilities', this.pathQuery(role, backend)).then((capabilities) => {
      if (!capabilities.canUpdate) {
        return this.router.transitionTo(ROUTES.VAULT_CLUSTER_SECRETS_BACKEND_LISTROOT, backend);
      }
      return this.store.createRecord('ssh-sign', {
        role: {
          backend,
          id: role,
          name: role,
        },
        id: `${backend}-${role}`,
      });
    });
  },

  setupController(controller) {
    this._super(...arguments);
    controller.set('backend', this.backendModel());
  },
});
