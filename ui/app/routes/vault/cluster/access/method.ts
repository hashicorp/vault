/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { supportedManagedAuthBackends } from 'vault/helpers/supported-managed-auth-backends';
import AuthMethodResource from 'vault/resources/auth/method';

import type ApiService from 'vault/services/api';
import type PathHelpService from 'vault/services/path-help';

export default class VaultClusterAccessMethodRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly pathHelp: PathHelpService;

  async model(params: { path: string }) {
    const { path } = params;
    const { auth } = await this.api.sys.internalUiListEnabledVisibleMounts();
    const methods = this.api
      .responseObjectToArray(auth, 'path')
      .map((method) => new AuthMethodResource(method, this));
    const method = methods.find((m) => m.id === path);
    // the user could have entered a random path in the URL that doesn't correspond to an existing method
    if (method) {
      const supportManaged = supportedManagedAuthBackends();
      // do not fetch path-help for unmanaged auth types
      if (!supportManaged.includes(method.methodType)) {
        method.paths = { apiPath: method.apiPath, paths: [], itemTypes: [] };
        return method;
      }
      return this.pathHelp.getPaths(method.apiPath, path, '', '').then((pathInfo) => {
        method.paths = pathInfo;
        return method;
      });
    } else {
      // throw a 404 if the path doesn't match any of the fetched methods
      throw { httpStatus: 404, path };
    }
  }
}
