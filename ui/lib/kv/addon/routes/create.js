/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
import { pathIsFromDirectory, breadcrumbsForDirectory } from 'vault/lib/kv-breadcrumbs';
import { withConfirmLeave } from 'core/decorators/confirm-leave';

@withConfirmLeave('model.secret')
export default class KvSecretsCreateRoute extends Route {
  @service store;
  @service secretMountPath;

  model(params) {
    const backend = this.secretMountPath.currentPath;
    const { initialKey: path } = params;

    return hash({
      backend,
      path,
      secret: this.store.createRecord('kv/data', { backend, path }),
      metadata: this.store.createRecord('kv/metadata', { backend, path }),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    let crumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'list' },
    ];

    if (pathIsFromDirectory(resolvedModel.path)) {
      crumbs = [...crumbs, ...breadcrumbsForDirectory(resolvedModel.path)];
    }
    crumbs.push({ label: 'create' });
    controller.breadcrumbs = crumbs;
  }
}
