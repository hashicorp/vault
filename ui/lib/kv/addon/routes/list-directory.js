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

  queryParams = {
    pageFilter: {
      refreshModel: true, // changing the "Filter secrets" input will cause the model hook to run again.
    },
  };

  model(params) {
    const pageFilter = params.pageFilter;
    const pathToSecret = params.path_to_secret ? normalizePath(params.path_to_secret) : '';
    const backend = this.secretMountPath.currentPath;
    const secrets = this.store
      .query('kv/metadata', { backend, pathToSecret })
      .then((models) => {
        this.has404 = false;
        // handle situation for when there is potentially both a pageFilter and pathToSecret ex: beep/my-.
        const filter = pathToSecret ? pathToSecret + (pageFilter || '') : pageFilter;
        return filter
          ? models.filter((model) => model.fullSecretPath.toLowerCase().includes(filter.toLowerCase()))
          : models;
      })
      .catch((err) => {
        if (err.httpStatus === 404) {
          this.has404 = true;
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
