/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class KvSecretsRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    // TODO: filter secrets based on queryParams value
    const backend = this.secretMountPath.get();
    const secrets = this.store.query('kv/metadata', { backend }).catch((error) => {
      if (error.httpStatus === 404) {
        return [];
      }
      throw error;
    });

    return hash({
      id: backend,
      backend,
      icon: 'kv',
      secrets,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend },
    ];
  }
}
