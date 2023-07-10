/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class KvSecretMetadataDiffRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const backend = this.secretMountPath.get();
    const parentModel = this.modelFor('secrets.secret.metadata');
    return hash({
      backend,
      ...parentModel,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'secrets' },
      { label: resolvedModel.name, route: 'secrets.secret.details', model: resolvedModel.name },
      { label: 'version diff' },
    ];
  }
}
