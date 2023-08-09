/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { pathIsFromDirectory, breadcrumbsForDirectory } from 'vault/lib/kv-breadcrumbs';

export default class KvSecretMetadataEditRoute extends Route {
  // model passed from 'secret' route, if we need to access or intercept
  // it can retrieved via `this.modelFor('secret'), which includes the metadata model.

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

    breadcrumbsArray = [
      ...breadcrumbsArray,
      { label: 'metadata', route: 'secret.details', model: resolvedModel.path },
      { label: 'edit' },
    ];
    controller.set('breadcrumbs', breadcrumbsArray);
  }
}
