/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
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
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'list', model: resolvedModel.backend },
      ...breadcrumbsForSecret(resolvedModel.backend, resolvedModel.path),
      { label: 'create' },
    ];
    controller.breadcrumbs = crumbs;
  }
}
