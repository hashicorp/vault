/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
import { normalizePath } from 'vault/utils/path-encoding-helpers';
import { breadcrumbsForDirectory } from 'vault/lib/kv-breadcrumbs';

export default class KvSecretsListRoute extends Route {
  @service store;
  @service router;
  @service secretMountPath;

  getPathToSecretFromUrl() {
    const { path_to_secret } = this.paramsFor('list-directory');
    return path_to_secret ? normalizePath(path_to_secret) : '';
  }

  model() {
    // TODO add filtering and return model for query on kv/metadata.
    const pathToSecret = this.getPathToSecretFromUrl();
    const backend = this.secretMountPath.currentPath;
    const arrayOfSecretModels = this.store.query('kv/metadata', { backend, pathToSecret }).catch((err) => {
      if (err.httpStatus === 404) {
        return [];
      } else {
        throw err;
      }
    });
    return hash({
      arrayOfSecretModels,
      backend,
      pathToSecret,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.routeName = this.routeName;
    let breadcrumbsArray = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'list' },
    ];
    // these breadcrumbs handle directories: beep/boop/
    if (resolvedModel.pathToSecret) {
      breadcrumbsArray = [...breadcrumbsArray, ...breadcrumbsForDirectory(resolvedModel.pathToSecret)];
    }
    controller.set('breadcrumbs', breadcrumbsArray);
  }
}
