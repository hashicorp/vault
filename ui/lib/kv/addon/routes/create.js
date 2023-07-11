/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class KvSecretsCreateRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    // TODO return model for query on kv/data
    const backend = this.secretMountPath.get();
    return hash({
      backend,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [{ label: resolvedModel.backend, route: 'list' }, { label: 'create' }];
  }
}
