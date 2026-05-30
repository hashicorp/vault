/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { IdentityApiOidcListProvidersListEnum } from '@hashicorp/vault-client-typescript';

export default class OidcProvidersRoute extends Route {
  @service api;
  @service capabilities;

  async model() {
    try {
      const response = await this.api.identity.oidcListProviders(IdentityApiOidcListProvidersListEnum.TRUE);
      const paths = response.keys.map((name) => this.capabilities.pathFor('oidcProvider', { name }));
      const capabilities = paths ? await this.capabilities.fetch(paths) : {};

      return {
        providers: this.api.keyInfoToArray(response, 'name'),
        capabilities,
      };
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        return {
          providers: [],
          capabilities: {},
        };
      } else {
        throw error;
      }
    }
  }
}
