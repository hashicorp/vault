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

  queryParams = {
    pageFilter: {
      refreshModel: true,
    },
  };

  async fetchMetadata(backend, pathToSecret, filter) {
    return await this.store
      .query('kv/metadata', { backend, pathToSecret })
      .then((models) => {
        return filter
          ? models.filter((model) => model.fullSecretPath.toLowerCase().includes(filter.toLowerCase()))
          : models;
      })
      .catch((err) => {
        if (err.httpStatus === 403) {
          return 403;
        }
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
  }

  model(params) {
    const pageFilter = params.pageFilter || '';
    const pathToSecret = params.path_to_secret ? normalizePath(params.path_to_secret) : '';
    const backend = this.secretMountPath.currentPath;
    const filter = pathToSecret ? pathToSecret + pageFilter : pageFilter;
    return hash({
      secrets: this.fetchMetadata(backend, pathToSecret, filter),
      backend,
      pathToSecret,
      filterValue: filter,
      pageFilter,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    if (resolvedModel.secrets === 403) {
      resolvedModel.noMetadataListPermissions = true;
    }
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

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.set('pageFilter', '');
    }
  }
}
