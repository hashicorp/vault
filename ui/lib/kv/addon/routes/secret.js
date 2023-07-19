/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class KvSecretRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const backend = this.secretMountPath.currentPath;
    const { name: path } = this.paramsFor('secret');
    return hash({
      path,
      backend,
      secret: this.store.queryRecord('kv/data', { backend, path }),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [{ label: resolvedModel.backend, route: 'list' }, { label: 'create' }];
  }
}
