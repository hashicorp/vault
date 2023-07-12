/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class KvSecretMetadataIndexRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    // TODO return model for query on kv/metadata.
    const backend = this.secretMountPath.get();
    const { name } = this.paramsFor('secret');
    return hash({
      path: name,
      backend,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'secret' },
      { label: resolvedModel.path, route: 'secret.details', model: resolvedModel.path },
      { label: 'metadata' },
    ];
  }
}
