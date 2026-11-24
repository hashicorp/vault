/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { IdentityApiOidcListScopesListEnum } from '@hashicorp/vault-client-typescript';

export default class OidcScopesRoute extends Route {
  @service api;
  @service capabilities;

  async model() {
    try {
      const { keys: scopes } = await this.api.identity.oidcListScopes(IdentityApiOidcListScopesListEnum.TRUE);
      const paths = scopes.map((name) => this.capabilities.pathFor('oidcScope', { name }));
      const capabilities = paths ? await this.capabilities.fetch(paths) : {};

      return {
        scopes,
        capabilities,
      };
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        return {
          scopes: [],
          capabilities: {},
        };
      } else {
        throw error;
      }
    }
  }
}
