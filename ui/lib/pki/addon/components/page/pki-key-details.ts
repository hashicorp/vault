/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';

import type RouterService from '@ember/routing/router-service';
import type FlashMessageService from 'vault/services/flash-messages';
import type { PkiReadKeyResponse } from '@hashicorp/vault-client-typescript';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';

interface Args {
  key: PkiReadKeyResponse;
}

export default class PkiKeyDetails extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  displayFields = ['key_id', 'key_name', 'key_type', 'key_bits'];

  get backend() {
    return this.secretMountPath.currentPath;
  }

  @action
  async deleteKey() {
    try {
      await this.api.secrets.pkiDeleteKey(this.args.key.key_id as string, this.backend);
      this.flashMessages.success('Key deleted successfully.');
      this.router.transitionTo('vault.cluster.secrets.backend.pki.keys.index');
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(message);
    }
  }
}
