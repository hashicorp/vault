/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { breadcrumbsForSecret } from 'kv/utils/kv-breadcrumbs';

export default class KvSecretPathsRoute extends Route {
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'list', model: resolvedModel.backend },
      ...breadcrumbsForSecret(resolvedModel.backend, resolvedModel.path),
      { label: 'Paths' },
    ];
  }
}
