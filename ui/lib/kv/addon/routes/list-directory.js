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

  model(params) {
    const pathToSecret = params.path_to_secret ? normalizePath(params.path_to_secret) : '';
    const backend = this.secretMountPath.currentPath;
    const secrets = this.store.query('kv/metadata', { backend, pathToSecret }).catch((err) => {
      if (err.httpStatus === 404) {
        return [];
      } else {
        throw err;
      }
    });
    return hash({
      secrets,
      backend,
      pathToSecret,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.routeName = this.routeName;

    let breadcrumbsArray = [{ label: 'secrets', route: 'secrets', linkExternal: true }];
    // if on top level don't link the engine breadcrumb label, but if within a directory, do link back to top level.
    if (this.routeName === 'list') {
      breadcrumbsArray.push({ label: resolvedModel.backend });
    } else {
      breadcrumbsArray.push({ label: resolvedModel.backend, route: 'list' });
    }
    // these breadcrumbs handle directories: beep/boop/
    if (resolvedModel.pathToSecret) {
      breadcrumbsArray = [...breadcrumbsArray, ...breadcrumbsForDirectory(resolvedModel.pathToSecret, true)];
    }
    controller.set('breadcrumbs', breadcrumbsArray);
    controller.set('routeName', this.routeName);
  }
}
