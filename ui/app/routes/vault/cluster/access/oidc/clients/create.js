/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import OidcClientForm from 'vault/forms/oidc/client';
import {
  IdentityApiOidcListAssignmentsListEnum,
  IdentityApiOidcListKeysListEnum,
} from '@hashicorp/vault-client-typescript';

export default class OidcClientsCreateRoute extends Route {
  @service api;

  async model() {
    // fetch keys and assignments to populate SearchSelect
    let keys = [];
    let assignments = [];
    try {
      const { keys: keyItems } = await this.api.identity.oidcListKeys(IdentityApiOidcListKeysListEnum.TRUE);
      const { keys: assignmentItems } = await this.api.identity.oidcListAssignments(
        IdentityApiOidcListAssignmentsListEnum.TRUE
      );
      // SearchSelect requires options to be objects
      keys = keyItems?.map((key) => ({ id: key }));
      assignments = assignmentItems
        ?.filter((assignment) => assignment !== 'allow_all')
        ?.map((assignment) => ({ id: assignment }));
    } catch (error) {
      // swallow error and return empty array for keys
    }
    const defaults = {
      key: 'default',
      id_token_ttl: '24h',
      access_token_ttl: '24h',
      client_type: 'confidential',
    };
    return {
      form: new OidcClientForm(defaults, { isNew: true }),
      keys,
      assignments,
    };
  }
}
