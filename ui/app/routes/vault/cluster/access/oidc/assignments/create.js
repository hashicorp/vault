/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import OidcAssignmentForm from 'vault/forms/oidc/assignment';

export default class OidcAssignmentsCreateRoute extends Route {
  model() {
    const { entities, groups } = this.modelFor('vault.cluster.access.oidc.assignments');
    return {
      form: new OidcAssignmentForm({}, { isNew: true }),
      entities,
      groups,
    };
  }
}
