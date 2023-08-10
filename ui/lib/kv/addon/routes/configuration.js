/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class KvConfigurationRoute extends Route {
  @service store;

  model() {
    const backend = this.modelFor('application');
    return hash({
      mountConfig: this.modelFor('application'),
      engineConfig: this.store.findRecord('kv/config', backend.id),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.mountConfig.id, route: 'list' },
      { label: 'configuration' },
    ];
  }
}
