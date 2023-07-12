/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
import { normalizePath } from 'vault/utils/path-encoding-helpers';

export default class KvSecretsListRoute extends Route {
  @service store;
  @service router;
  @service secretMountPath;

  getSecretPrefixFromUrl() {
    const { secret_prefix } = this.paramsFor('list-nested-secret');
    return secret_prefix ? normalizePath(secret_prefix) : '';
  }

  model() {
    // TODO add filtering and return model for query on kv/metadata.
    const secretPrefix = this.getSecretPrefixFromUrl();
    const backend = this.secretMountPath.currentPath;
    const arrayOfSecretModels = this.store.query('kv/metadata', { backend, secretPrefix }).catch((err) => {
      if (err.httpStatus === 404) {
        return [];
      } else {
        throw err;
      }
    });
    return hash({
      arrayOfSecretModels,
      backend,
      routeName: this.routeName,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.set('model', resolvedModel.arrayOfSecretModels);
    controller.routeName = resolvedModel.routeName;
    controller.pageTitle = resolvedModel.backend;
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend },
    ];
  }
}
