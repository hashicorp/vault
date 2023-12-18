/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { breadcrumbsForSecret } from 'kv/utils/kv-breadcrumbs';

export default class KvSecretMetadataEditRoute extends Route {
  // model passed from 'secret' route, if we need to access or intercept
  // it can retrieved via `this.modelFor('secret'), which includes the metadata model.

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const breadcrumbsArray = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'list' },
      ...breadcrumbsForSecret(resolvedModel.path),
      { label: 'metadata', route: 'secret.metadata' },
      { label: 'edit' },
    ];

    controller.set('breadcrumbs', breadcrumbsArray);
  }
}
