/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class KvConfigurationRoute extends Route {
  @service api;

  async model() {
    const backend = this.modelFor('application');
    // display mount config if engine config request fails
    const engineConfig = await this.api.secrets.kvV2ReadConfiguration(backend.id).catch(() => {});

    return {
      ...engineConfig,
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const backend = this.modelFor('application');
    controller.backend = backend;
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: backend.id, route: 'list', model: backend.id },
      { label: 'Configuration' },
    ];
  }
}
