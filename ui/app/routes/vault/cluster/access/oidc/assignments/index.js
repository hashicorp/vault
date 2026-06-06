/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { IdentityApiOidcListAssignmentsListEnum } from '@hashicorp/vault-client-typescript';

export default class OidcAssignmentsIndexRoute extends Route {
  @service api;
  @service capabilities;

  async model() {
    try {
      const { keys: assignments } = await this.api.identity.oidcListAssignments(
        IdentityApiOidcListAssignmentsListEnum.TRUE
      );
      const paths = assignments.map((name) => this.capabilities.pathFor('oidcAssignment', { name }));
      const capabilities = paths ? await this.capabilities.fetch(paths) : {};

      return {
        assignments,
        capabilities,
      };
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        return {
          assignments: [],
          capabilities: {},
        };
      } else {
        throw error;
      }
    }
  }
}
