/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { breadcrumbsForSecret } from 'kv/utils/kv-breadcrumbs';
import { service } from '@ember/service';
import { ROUTES } from 'vault/utils/routes';

export default class SecretPatch extends Route {
  @service('app-router') router;

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const breadcrumbsArray = [
      { label: 'Secrets', route: ROUTES.SECRETS, linkExternal: true },
      { label: resolvedModel.backend, route: ROUTES.LIST, model: resolvedModel.backend },
      ...breadcrumbsForSecret(resolvedModel.backend, resolvedModel.path),
      { label: 'Patch' },
    ];
    controller.breadcrumbs = breadcrumbsArray;
  }

  // isPatchAllowed is true if (1) the version is enterprise, (2) a user has "patch" secret + "read" subkeys capabilities, (3) latest secret version is not deleted or destroyed
  redirect(model) {
    if (!model.isPatchAllowed) {
      this.router.transitionTo(ROUTES.VAULT_CLUSTER_SECRETS_BACKEND_KV_SECRET_INDEX, model.path);
    }
  }
}
