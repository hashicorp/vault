/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import AdapterError from '@ember-data/adapter/error';
import { breadcrumbsForSecret } from 'kv/utils/kv-breadcrumbs';
import { service } from '@ember/service';
export default class SecretPatch extends Route {
  @service version;

  beforeModel() {
    if (!this.version.isEnterprise) {
      throw this.generateError({ message: 'Patching a KV v2 secret is only available on Vault Enterprise.' });
    }
  }

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

  afterModel(model) {
    if (model.canPatchSecret && !!model.subkeys.subkeys) {
      return;
    }
    throw this.generateError({
      title: 'You do not have permissions to patch a KV v2 secret',
      message: 'Ask your administrator if you think you should have access to:',
      permissionsError: {
        READ: `"${model.backend}/subkeys/${model.path}"`,
        PATCH: `"${model.backend}/data/${model.path}"`,
      },
    });
  }

  generateError({ title, message, permissionsError }) {
    const error = new AdapterError();
    if (title) error.title = title;
    if (message) error.message = message;
    if (permissionsError) error.permissionsError = permissionsError;
    // since we customized this error we don't need the built-in Ember errors
    delete error.errors;
    return error;
  }
}
