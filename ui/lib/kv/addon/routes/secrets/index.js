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
    // TODO add filtering and return model for query on kv/metadata.
    const backend = this.secretMountPath.get();
    return hash({
      id: backend,
      backend,
      icon: 'kv',
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
