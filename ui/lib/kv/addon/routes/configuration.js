/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class KvConfigurationRoute extends Route {
  @service store;

  model() {
    const engine = this.modelFor('application');
    return this.store
      .query('secret-engine', {
        path: engine.id,
      })
      .then((engine) => {
        if (engine) {
          return engine.get('firstObject');
        }
      });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.pageTitle = resolvedModel.backend;
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'list' },
      { label: 'configuration' },
    ];
  }
}
