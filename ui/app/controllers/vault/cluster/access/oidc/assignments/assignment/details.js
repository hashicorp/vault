/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { action } from '@ember/object';
import { service } from '@ember/service';

export default class OidcAssignmentDetailsController extends Controller {
  @service api;
  @service router;
  @service flashMessages;

  @action
  async delete() {
    try {
      await this.api.identity.oidcDeleteAssignment(this.model.assignment.name);
      this.flashMessages.success('Assignment deleted successfully');
      this.router.transitionTo('vault.cluster.access.oidc.assignments');
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(message);
    }
  }
}
