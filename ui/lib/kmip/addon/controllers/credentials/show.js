/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';

export default class CredentialsShowController extends Controller {
  @service flashMessages;
  @service router;

  @action
  async revokeCredentials() {
    try {
      await this.model.destroyRecord();
      this.flashMessages.success('Successfully revoked credentials');
      this.router.transitionTo('vault.cluster.secrets.backend.kmip.credentials.index', this.scope, this.role);
    } catch (e) {
      this.flashMessages.danger(`There was an error revoking credentials: ${e.errors.join(' ')}`);
      this.model.rollbackAttributes();
    }
  }
}
