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
    const { name } = this.paramsFor('secrets.secret');
    return hash({
      id: name,
      backend,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'secrets' },
      { label: resolvedModel.id, route: 'secrets.secret.details', model: resolvedModel.id },
      { label: 'metadata' },
    ];
  }
}
