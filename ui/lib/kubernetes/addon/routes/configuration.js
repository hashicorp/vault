/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import FetchConfigRoute from './fetch-config';

export default class KubernetesConfigureRoute extends FetchConfigRoute {
  model() {
    return {
      backend: this.modelFor('application'),
      config: this.configModel,
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend.id },
    ];
  }
}
