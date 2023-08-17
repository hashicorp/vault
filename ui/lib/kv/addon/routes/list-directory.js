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
    currentPage: {
      refreshModel: true,
    },
  };

  async fetchMetadata(backend, pathToSecret, params) {
    return await this.store
      .lazyPaginatedQuery('kv/metadata', {
        backend,
        responsePath: 'data.keys',
        page: params.currentPage || 1,
        pageFilter: params.pageFilter,
        pathToSecret,
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
    let { pageFilter, path_to_secret } = params;
    pageFilter = pageFilter || '';
    const pathToSecret = path_to_secret ? normalizePath(path_to_secret) : '';
    const backend = this.secretMountPath.currentPath;
    const filter = pathToSecret ? pathToSecret + pageFilter : pageFilter;
    return hash({
      secrets: this.fetchMetadata(backend, pathToSecret, params),
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
      controller.set('pageFilter', null);
      controller.set('currentPage', null);
    }
  }
}
