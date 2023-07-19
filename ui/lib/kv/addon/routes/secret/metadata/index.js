/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
import { pathIsFromDirectory, breadcrumbsForDirectory } from 'vault/lib/kv-breadcrumbs';

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

    let breadcrumbsArray = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'list' },
    ];

    if (pathIsFromDirectory(resolvedModel.path)) {
      breadcrumbsArray = [...breadcrumbsArray, ...breadcrumbsForDirectory(resolvedModel.path)];
    } else {
      breadcrumbsArray.push({
        label: resolvedModel.path,
        route: 'secret.details',
        model: resolvedModel.path,
      });
    }

    breadcrumbsArray.push({ label: 'metadata' });
    controller.set('breadcrumbs', breadcrumbsArray);
  }
}
