/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class KvSecretsCreateRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const backend = this.secretMountPath.get();
    return this.store.createRecord('kv/secret', { backend });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    // ARG need to link in first row to vault.cluster.secrets.backends
    controller.breadcrumbs = [
      // { label: resolvedModel.backend, route: 'backends' },
      { label: 'secrets', route: 'secrets' },
      { label: 'create' },
    ];
  }
}
