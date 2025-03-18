/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';
import { ROUTES } from 'vault/utils/routes';

export default class KvConfigurationRoute extends Route {
  @service store;

  model() {
    const backend = this.modelFor('application');
    return hash({
      mountConfig: backend,
      engineConfig: this.store.findRecord('kv/config', backend.id).catch(() => {
        // return an empty record so we have access to model capabilities
        return this.store.createRecord('kv/config', { backend: backend.id });
      }),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Secrets', route: ROUTES.SECRETS, linkExternal: true },
      { label: resolvedModel.mountConfig.id, route: ROUTES.LIST, model: resolvedModel.engineConfig.backend },
      { label: 'Configuration' },
    ];
  }
}
