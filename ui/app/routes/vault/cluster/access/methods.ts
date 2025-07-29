/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import AuthMethodResource from 'vault/resources/auth/method';

import type ApiService from 'vault/services/api';
import type Capabilities from 'vault/services/capabilities';

export default class VaultClusterAccessMethodsRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly capabilities: Capabilities;

  queryParams = {
    page: {
      refreshModel: true,
    },
    pageFilter: {
      refreshModel: true,
    },
  };

  async model() {
    const { auth } = await this.api.sys.internalUiListEnabledVisibleMounts();

    const methods = this.api
      .responseObjectToArray(auth, 'path')
      .map((method) => new AuthMethodResource(method, this));

    const paths = methods.reduce((paths: string[], { path, methodType }) => {
      paths.push(
        this.capabilities.pathFor('authMethodConfig', { path }),
        this.capabilities.pathFor('authMethodDelete', { path })
      );
      if (methodType === 'aws') {
        paths.push(this.capabilities.pathFor('authMethodConfigAws', { path }));
      }
      return paths;
    }, []);

    const capabilities = this.capabilities.fetch(paths);

    return { methods, capabilities };
  }
}
