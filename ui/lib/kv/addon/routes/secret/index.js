/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { breadcrumbsForSecret } from 'kv/utils/kv-breadcrumbs';
import { ROUTES } from 'vault/utils/routes';

export default class SecretIndex extends Route {
  @service('app-router') router;

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const breadcrumbsArray = [
      { label: 'Secrets', route: ROUTES.SECRETS, linkExternal: true },
      { label: resolvedModel.backend, route: ROUTES.LIST, model: resolvedModel.backend },
      ...breadcrumbsForSecret(resolvedModel.backend, resolvedModel.path, true),
    ];
    controller.breadcrumbs = breadcrumbsArray;
  }
}
