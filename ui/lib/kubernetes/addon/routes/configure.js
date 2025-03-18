/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { withConfig } from 'core/decorators/fetch-secrets-engine-config';
import { ROUTES } from 'vault/utils/routes';

@withConfig('kubernetes/config')
export default class KubernetesConfigureRoute extends Route {
  @service store;
  @service secretMountPath;

  async model() {
    const backend = this.secretMountPath.currentPath;
    return this.configModel || this.store.createRecord('kubernetes/config', { backend });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Secrets', route: ROUTES.SECRETS, linkExternal: true },
      { label: resolvedModel.backend, route: ROUTES.OVERVIEW, model: resolvedModel.backend },
      { label: 'Configure' },
    ];
  }
}
