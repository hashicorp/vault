/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import OidcAssignmentForm from 'vault/forms/oidc/assignment';

export default class OidcAssignmentEditRoute extends Route {
  model() {
    const { entities, groups } = this.modelFor('vault.cluster.access.oidc.assignments');
    const { assignment } = this.modelFor('vault.cluster.access.oidc.assignments.assignment');
    return {
      form: new OidcAssignmentForm(assignment),
      entities,
      groups,
    };
  }
}
