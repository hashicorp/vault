/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { breadcrumbsForSecret } from 'kv/utils/kv-breadcrumbs';
import { service } from '@ember/service';
export default class SecretPatch extends Route {
  @service router;
  @service version;

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const breadcrumbsArray = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'list', model: resolvedModel.backend },
      ...breadcrumbsForSecret(resolvedModel.backend, resolvedModel.path),
      { label: 'Patch' },
    ];
    controller.breadcrumbs = breadcrumbsArray;
  }

  // patch is enterprise only and users must permissions to "read" subkeys and "patch" a secret
  redirect(model) {
    if (!this.version.isEnterprise || !model.canPatchSecret || !model.subkeys.subkeys) {
      this.router.transitionTo('vault.cluster.secrets.backend.kv.secret.index', model.path);
    }
  }
}
