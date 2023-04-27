/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfig } from 'kubernetes/decorators/fetch-config';
import { hash } from 'rsvp';

@withConfig()
export default class KubernetesRolesRoute extends Route {
  @service store;
  @service secretMountPath;

  model(params, transition) {
    // filter roles based on pageFilter value
    const { pageFilter } = transition.to.queryParams;
    const roles = this.store
      .query('kubernetes/role', { backend: this.secretMountPath.get() })
      .then((models) =>
        pageFilter
          ? models.filter((model) => model.name.toLowerCase().includes(pageFilter.toLowerCase()))
          : models
      )
      .catch((error) => {
        if (error.httpStatus === 404) {
          return [];
        }
        throw error;
      });
    return hash({
      backend: this.modelFor('application'),
      promptConfig: this.promptConfig,
      roles,
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
