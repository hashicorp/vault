/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfig } from 'core/decorators/fetch-secrets-engine-config';
import { hash } from 'rsvp';

@withConfig('kubernetes/config')
export default class KubernetesOverviewRoute extends Route {
  @service store;
  @service secretMountPath;

  async model() {
    const backend = this.secretMountPath.currentPath;
    return hash({
      promptConfig: this.promptConfig,
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
