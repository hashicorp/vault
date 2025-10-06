/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { breadcrumbsForSecret } from 'kv/utils/kv-breadcrumbs';
import KvForm from 'vault/forms/secrets/kv';

export default class KvSecretDetailsEditRoute extends Route {
  model() {
    const parentModel = this.modelFor('secret.details');
    const { metadata, secret } = parentModel;
    const formData = {
      path: parentModel.path,
      max_versions: 0,
      options: {
        cas: metadata?.current_version || secret.version,
      },
    };
    if (!parentModel.secret.failReadErrorCode) {
      formData.secretData = parentModel.secret.secretData;
    }
    return {
      ...parentModel,
      form: new KvForm(formData),
    };
  }
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'list', model: resolvedModel.backend },
      ...breadcrumbsForSecret(resolvedModel.backend, resolvedModel.path),
      { label: 'Edit' },
    ];
  }
}
