/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
import { withConfirmLeave } from 'core/decorators/confirm-leave';
import { breadcrumbsForSecret } from 'kv/utils/kv-breadcrumbs';

@withConfirmLeave('model.secret', ['model.metadata'])
export default class KvSecretsCreateRoute extends Route {
  @service store;
  @service secretMountPath;

  model(params) {
    const backend = this.secretMountPath.currentPath;
    const { initialKey: path } = params;

    return hash({
      backend,
      path,
      // see serializer for logic behind setting casVersion
      secret: this.store.createRecord('kv/data', { backend, path, casVersion: 0 }),
      metadata: this.store.createRecord('kv/metadata', { backend, path }),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    const crumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'list' },
      ...breadcrumbsForSecret(resolvedModel.path),
      { label: 'create' },
    ];
    controller.breadcrumbs = crumbs;
  }
}
