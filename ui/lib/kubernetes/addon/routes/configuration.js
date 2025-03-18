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

  model() {
    // in case of any error other than 404 we want to display that to the user
    if (this.configError) {
      throw this.configError;
    }
    return {
      backend: this.modelFor('application'),
      config: this.configModel,
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Secrets', route: ROUTES.SECRETS, linkExternal: true },
      { label: resolvedModel.backend.id, route: ROUTES.OVERVIEW, model: resolvedModel.backend },
      { label: 'Configuration' },
    ];
  }
}
