/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';

export default class KmipPageCredentialsGenerate extends Component {
  @service flashMessages;
  @service('app-router') router;
  @tracked hasGenerated = false;

  @action
  async revokeCredentials() {
    const { scope, role } = this.args.credentials;
    try {
      await this.args.credentials.destroyRecord();
      this.flashMessages.success('Successfully revoked credentials.');
      this.router.transitionTo('vault.cluster.secrets.backend.kmip.credentials.index', scope, role);
    } catch (e) {
      const message = errorMessage(e);
      this.flashMessages.danger(`There was an error revoking credentials: ${message}.`);
      this.args.credentials.rollbackAttributes();
    }
  }
}
