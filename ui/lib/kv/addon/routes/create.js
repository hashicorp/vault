/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { breadcrumbsForSecret } from 'kv/utils/kv-breadcrumbs';
import KvForm from 'vault/forms/secrets/kv';

export default class KvSecretsCreateRoute extends Route {
  @service secretMountPath;

  model(params) {
    const backend = this.secretMountPath.currentPath;
    const { initialKey: path } = params;

    return {
      backend,
      path,
      form: new KvForm(
        {
          path,
          max_versions: 0,
          delete_version_after: '0s',
          cas_required: false,
          options: { cas: 0 },
        },
        { isNew: true }
      ),
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    const crumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'list', model: resolvedModel.backend },
      ...breadcrumbsForSecret(resolvedModel.backend, resolvedModel.path),
      { label: 'Create' },
    ];
    controller.breadcrumbs = crumbs;
  }
}
