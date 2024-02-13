/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';
import type RouterService from '@ember/routing/router-service';
import type FlashMessageService from 'vault/services/flash-messages';
import type PkiKeyModel from 'vault/models/pki/key';
interface Args {
  key: PkiKeyModel;
}

export default class PkiKeyDetails extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;

  @action
  async deleteKey() {
    try {
      await this.args.key.destroyRecord();
      this.flashMessages.success('Key deleted successfully.');
      this.router.transitionTo('vault.cluster.secrets.backend.pki.keys.index');
    } catch (error) {
      this.flashMessages.danger(errorMessage(error));
    }
  }
}
