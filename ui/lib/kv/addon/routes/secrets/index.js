/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class KvSecretsRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    // TODO add filtering and return model for query on kv/metadata.
    const backend = this.secretMountPath.currentPath;
    const secrets = this.store.query('kv/metadata', { backend }).catch((err) => {
      if (err.httpStatus === 404) {
        return [];
      } else {
        throw err;
      }
    });
    return hash({
      secrets,
      backend,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.set('model', resolvedModel.secrets);
    controller.pageTitle = resolvedModel.backend;
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend },
    ];
  }
}
