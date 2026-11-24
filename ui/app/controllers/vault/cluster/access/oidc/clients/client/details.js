/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { action } from '@ember/object';
import { service } from '@ember/service';

export default class OidcClientDetailsController extends Controller {
  @service api;
  @service router;
  @service flashMessages;

  @action
  async delete() {
    try {
      await this.api.identity.oidcDeleteClient(this.model.client.name);
      this.flashMessages.success('Application deleted successfully');
      this.router.transitionTo('vault.cluster.access.oidc.clients');
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(message);
    }
  }
}
