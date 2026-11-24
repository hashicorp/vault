/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { IdentityApiOidcListKeysListEnum } from '@hashicorp/vault-client-typescript';

export default class OidcKeysRoute extends Route {
  @service api;
  @service capabilities;

  async model() {
    try {
      const { keys } = await this.api.identity.oidcListKeys(IdentityApiOidcListKeysListEnum.TRUE);
      const paths = keys.map((name) => this.capabilities.pathFor('oidcKey', { name }));
      const capabilities = paths ? await this.capabilities.fetch(paths) : {};

      return {
        keys,
        capabilities,
      };
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        return {
          keys: [],
          capabilities: {},
        };
      } else {
        throw error;
      }
    }
  }
}
