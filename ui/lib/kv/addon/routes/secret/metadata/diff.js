/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { breadcrumbsForSecret } from 'kv/utils/kv-breadcrumbs';

export default class KvSecretMetadataDiffRoute extends Route {
  // model passed from parent secret route, if we need to access or intercept
  // it can retrieved via `this.modelFor('secret'), which includes the metadata model.
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const breadcrumbsArray = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'list', model: resolvedModel.backend },
      ...breadcrumbsForSecret(resolvedModel.backend, resolvedModel.path),
      { label: 'version history', route: 'secret.metadata.versions' },
      { label: 'diff' },
    ];
    controller.set('breadcrumbs', breadcrumbsArray);
  }
}
