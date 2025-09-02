/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';
import { pathIsDirectory, breadcrumbsForSecret } from 'kv/utils/kv-breadcrumbs';

export default class KvSecretsListRoute extends Route {
  @service pagination;
  @service('app-router') router;
  @service secretMountPath;

  queryParams = {
    pageFilter: {
      refreshModel: true,
    },
    page: {
      refreshModel: true,
    },
  };

  async fetchMetadata(backend, pathToSecret, params) {
    return await this.pagination
      .lazyPaginatedQuery('kv/metadata', {
        backend,
        responsePath: 'data.keys',
        page: Number(params.page) || 1,
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

  getPathToSecret(pathParam) {
    if (!pathParam) return '';
    // links and routing assumes pathToParam includes trailing slash
    // users may want to include a percent-encoded octet like %2f in their path. Example: 'foo%2fbar' or non-data octets like 'foo%bar'.
    // we are assuming the user intended to include these characters in their path and we should not decode them.
    return pathIsDirectory(pathParam) ? pathParam : `${pathParam}/`;
  }

  model(params) {
    const { pageFilter, path_to_secret } = params;
    const pathToSecret = this.getPathToSecret(path_to_secret);
    const backend = this.secretMountPath.currentPath;
    const filterValue = pathToSecret ? (pageFilter ? pathToSecret + pageFilter : pathToSecret) : pageFilter;
    return hash({
      secrets: this.fetchMetadata(backend, pathToSecret, params),
      backend,
      pathToSecret,
      filterValue,
      pageFilter,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    // renders alert inline error for overview card
    resolvedModel.failedDirectoryQuery =
      resolvedModel.secrets === 403 && pathIsDirectory(resolvedModel.pathToSecret);

    let breadcrumbsArray = [{ label: 'Secrets', route: 'secrets', linkExternal: true }];
    // if on top level don't link the engine breadcrumb label, but if within a directory, do link back to top level.
    if (this.routeName === 'list') {
      breadcrumbsArray.push({ label: resolvedModel.backend });
    } else {
      breadcrumbsArray = [
        ...breadcrumbsArray,
        { label: resolvedModel.backend, route: 'list', model: resolvedModel.backend },
        ...breadcrumbsForSecret(resolvedModel.backend, resolvedModel.pathToSecret, true),
      ];
    }

    controller.set('breadcrumbs', breadcrumbsArray);
  }

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.set('pageFilter', null);
      controller.set('page', null);
    }
  }
}
