/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { breadcrumbsForSecret } from 'kv/utils/kv-breadcrumbs';

export default class KvSecretMetadataEditRoute extends Route {
  // model passed from 'secret' route, if we need to access or intercept
  // it can retrieved via `this.modelFor('secret'), which includes the metadata model.

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const breadcrumbsArray = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'list', model: resolvedModel.backend },
      ...breadcrumbsForSecret(resolvedModel.backend, resolvedModel.path),
      { label: 'Metadata', route: 'secret.metadata' },
      { label: 'Edit' },
    ];

    controller.set('breadcrumbs', breadcrumbsArray);
  }
}
