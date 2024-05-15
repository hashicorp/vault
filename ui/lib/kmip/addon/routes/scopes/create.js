/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class KmipScopesCreate extends Route {
  @service store;
  @service secretMountPath;

  beforeModel() {
    this.store.unloadAll('kmip/scope');
  }

  model() {
    const backend = this.secretMountPath.currentPath;
    return this.store.createRecord('kmip/scope', {
      backend: backend,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    const crumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'scopes', model: resolvedModel.backend },
      { label: 'create' },
    ];
    controller.breadcrumbs = crumbs;
  }
}
