/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class OidcAssignmentRoute extends Route {
  @service api;
  @service capabilities;

  async model({ name }) {
    const { data } = await this.api.identity.oidcReadAssignment(name);
    const capabilities = await this.capabilities.for('oidcAssignment', { name });
    return {
      assignment: { ...data, name },
      capabilities,
    };
  }
}
