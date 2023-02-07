/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { hash } from 'rsvp';
import FetchConfigRoute from './fetch-config';

export default class KubernetesOverviewRoute extends FetchConfigRoute {
  async model() {
    const backend = this.secretMountPath.get();
    return hash({
      config: this.configModel,
      backend: this.modelFor('application'),
      roles: this.store.query('kubernetes/role', { backend }).catch(() => []),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend.id },
    ];
  }
}
