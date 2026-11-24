/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { hash } from 'rsvp';
import { service } from '@ember/service';
import { getEnginePathParam } from 'vault/utils/backend-route-helpers';
import {
  SecretsApiDatabaseListStaticRolesListEnum,
  SecretsApiDatabaseListRolesListEnum,
  SecretsApiDatabaseListConnectionsListEnum,
} from '@hashicorp/vault-client-typescript';

export default Route.extend({
  capabilities: service(),
  api: service(),
  type: '',

  // this only grabs connections for current db backend, only used to populate # of connections
  async fetchConnection(queryOptions) {
    try {
      const { keys } = await this.api.secrets.databaseListConnections(
        queryOptions.backend,
        SecretsApiDatabaseListConnectionsListEnum.TRUE
      );
      return keys;
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        return status;
      }
    }
  },

  // this grabs both dynamic and static roles for current db backend, only used to populate # of roles
  async fetchAllRoles(queryOptions) {
    try {
      const roles = [];
      const { backend } = queryOptions;
      const [staticResp, dynamicResp] = await Promise.allSettled([
        this.api.secrets.databaseListStaticRoles(backend, SecretsApiDatabaseListStaticRolesListEnum.TRUE),
        this.api.secrets.databaseListRoles(backend, SecretsApiDatabaseListRolesListEnum.TRUE),
      ]);

      if (staticResp.status === 'rejected' && dynamicResp.status === 'rejected') {
        const { response: staticError, status: staticStatus } = await this.api.parseError(staticResp.reason);
        const { response: dynamicError, status: dynamicStatus } = await this.api.parseError(
          dynamicResp.reason
        );
        if (staticError?.isControlGroupError) {
          throw staticError;
        }
        throw staticStatus < dynamicStatus ? dynamicError : staticError;
      } else {
        if (staticResp.value) {
          roles.push(...staticResp.value.keys);
        }
        if (dynamicResp.value) {
          roles.push(...dynamicResp.value.keys);
        }
        return roles;
      }
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        return status;
      }
    }
  },

  async fetchCapabilitiesRole(queryOptions) {
    const paths = [this.capabilities.pathFor('databaseRoles', { backend: queryOptions.backend })];
    const capabilities = paths ? await this.capabilities.fetch(paths) : {};
    return capabilities[paths[0]];
  },

  async fetchCapabilitiesStaticRole(queryOptions) {
    const paths = [this.capabilities.pathFor('databaseStaticRoles', { backend: queryOptions.backend })];
    const capabilities = paths ? await this.capabilities.fetch(paths) : {};
    return capabilities[paths[0]];
  },

  async fetchCapabilitiesConnection(queryOptions) {
    const paths = [this.capabilities.pathFor('databaseConfig', { backend: queryOptions.backend })];
    const capabilities = paths ? await this.capabilities.fetch(paths) : {};
    return capabilities[paths[0]];
  },

  async model() {
    const backend = getEnginePathParam(this);
    const queryOptions = { backend, id: '' };

    const connection = await this.fetchConnection(queryOptions);
    const role = await this.fetchAllRoles(queryOptions);
    const roleCapabilities = await this.fetchCapabilitiesRole(queryOptions);
    const staticRoleCapabilities = await this.fetchCapabilitiesStaticRole(queryOptions);
    const connectionCapabilities = await this.fetchCapabilitiesConnection(queryOptions);

    return hash({
      backend,
      connections: connection,
      roles: role,
      engineType: 'database',
      id: backend,
      roleCapabilities,
      staticRoleCapabilities,
      connectionCapabilities,
      icon: 'database',
    });
  },

  setupController(controller, model) {
    this._super(...arguments);
    const showEmptyState = model.connections === 404 && (model.roles === undefined || model.roles === 404);
    const noConnectionCapabilities =
      !model.connectionCapabilities.canList &&
      !model.connectionCapabilities.canCreate &&
      !model.connectionCapabilities.canUpdate;

    const emptyStateMessage = function () {
      if (noConnectionCapabilities) {
        return 'You cannot yet generate credentials.  Ask your administrator if you think you should have access.';
      } else {
        return 'You can connect an external database to Vault.  We recommend that you create a user for Vault rather than using the database root user.';
      }
    };
    controller.set('showEmptyState', showEmptyState);
    controller.set('emptyStateMessage', emptyStateMessage());
  },
});
