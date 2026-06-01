/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default Route.extend({
  router: service(),
  capabilities: service(),
  templateName: 'vault/cluster/secrets/backend/sign',

  backendModel() {
    return this.modelFor('vault.cluster.secrets.backend');
  },

  async model(params) {
    const role = params.secret;
    const backendModel = this.backendModel();
    const backend = backendModel.id;

    if (backendModel.type !== 'ssh') {
      return this.router.transitionTo('vault.cluster.secrets.backend.list-root', backend);
    }

    const signPath = this.capabilities.pathFor('sshSign', { backend, id: role });
    const capabilities = await this.capabilities.fetch([signPath]);
    if (!capabilities[signPath]?.canUpdate) {
      return this.router.transitionTo('vault.cluster.secrets.backend.list-root', backend);
    }

    return {
      roleName: role,
      backendPath: backend,
    };
  },
});
