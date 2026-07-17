/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class OidcAssignmentDetailsRoute extends Route {
  @service api;

  async model() {
    const { entities, groups } = this.modelFor('vault.cluster.access.oidc.assignments');
    const { assignment, capabilities } = this.modelFor('vault.cluster.access.oidc.assignments.assignment');
    assignment.entities = assignment.entity_ids.map(
      (id) => entities.find((entity) => entity.id === id)?.name ?? id
    );
    assignment.groups = assignment.group_ids.map((id) => groups.find((group) => group.id === id)?.name ?? id);

    return {
      assignment,
      capabilities,
    };
  }
}
