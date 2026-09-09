/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { action } from '@ember/object';
import { service } from '@ember/service';

export default class OidcScopeDetailsController extends Controller {
  @service api;
  @service router;
  @service flashMessages;

  @action
  async delete() {
    try {
      await this.api.identity.oidcDeleteScope(this.model.scope.name);
      this.flashMessages.success('Scope deleted successfully');
      this.router.transitionTo('vault.cluster.access.oidc.scopes');
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(message);
    }
  }
}
