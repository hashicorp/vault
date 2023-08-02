/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
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
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [{ label: resolvedModel.backend, route: 'list' }, { label: 'create' }];
  }
}
