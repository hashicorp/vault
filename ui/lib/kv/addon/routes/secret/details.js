/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
import { pathIsFromNested, modelIdsForSecretPrefix } from 'vault/lib/kv-breadcrumbs';

export default class KvSecretDetailsRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    // TODO return model for query on kv/data.
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

    if (pathIsFromNested(resolvedModel.path)) {
      const breadcrumbsForSecretPrefix = modelIdsForSecretPrefix(resolvedModel.path);
      breadcrumbsArray = [...breadcrumbsArray, ...breadcrumbsForSecretPrefix];
    } else {
      const breadcrumbsCurrentPath = { label: resolvedModel.path };
      breadcrumbsArray.push(breadcrumbsCurrentPath);
    }
    controller.breadcrumbs = breadcrumbsArray;
  }
}
