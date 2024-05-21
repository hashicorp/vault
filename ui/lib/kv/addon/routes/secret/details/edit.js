/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';
import { withConfirmLeave } from 'core/decorators/confirm-leave';
import { breadcrumbsForSecret } from 'kv/utils/kv-breadcrumbs';

@withConfirmLeave('model.newVersion')
export default class KvSecretDetailsEditRoute extends Route {
  @service store;

  model() {
    const parentModel = this.modelFor('secret.details');
    const { backend, path, secret, metadata } = parentModel;
    return hash({
      secret,
      metadata,
      backend,
      path,
      newVersion: this.store.createRecord('kv/data', {
        backend,
        path,
        secretData: secret?.secretData,
        // see serializer for logic behind setting casVersion
        casVersion: metadata?.currentVersion || secret?.version,
      }),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'list', model: resolvedModel.backend },
      ...breadcrumbsForSecret(resolvedModel.backend, resolvedModel.path),
      { label: 'edit' },
    ];
  }
}
