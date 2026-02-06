/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import KvForm from 'vault/forms/secrets/kv';
import { service } from '@ember/service';

export default class KvConfigureRoute extends Route {
  @service api;
  @service('app-router') router;

  async model() {
    const backend = this.modelFor('application');
    const engineConfig = await this.api.secrets.kvV2ReadConfiguration(backend.id).catch(() => {});
    const { max_versions, cas_required, delete_version_after } = engineConfig;

    return {
      form: new KvForm({ path: backend.id, max_versions, cas_required, delete_version_after }),
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const backend = this.modelFor('application');
    controller.backend = backend;
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: backend.id, route: 'list', model: backend.id },
      { label: 'Configuration', route: 'configuration', model: backend },
      { label: 'Edit' },
    ];
  }
}
