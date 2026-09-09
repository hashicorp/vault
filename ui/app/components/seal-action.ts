/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';

import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';

interface Args {
  onSeal: CallableFunction;
}

export default class SealActionComponent extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @action
  async handleSeal() {
    try {
      await this.args.onSeal();
    } catch (error) {
      const message = await this.api.parseError(error, 'Check Vault logs for details.');

      this.flashMessages.danger(message.message, {
        title: 'Seal attempt failed',
      });
    }
  }
}
