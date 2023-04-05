/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class PkiConfigurationTidyRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    return hash({
      engine: this.modelFor('application'),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'configuration', route: 'configuration.index' },
      { label: 'tidy' },
    ];
  }
}
