/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

import type RouterService from '@ember/routing/router-service';
import type FlashMessageService from 'vault/services/flash-messages';
import type ApiService from 'vault/services/api';
import type VersionService from 'vault/services/version';

interface Args {
  backend: string;
}

export default class PkiConfigurationDetails extends Component<Args> {
  @service declare readonly api: ApiService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly version: VersionService;

  @tracked showDeleteAllIssuers = false;

  get isEnterprise() {
    return this.version.isEnterprise;
  }

  @action
  async deleteAllIssuers() {
    try {
      await this.api.secrets.pkiDeleteRoot(this.args.backend);
      this.flashMessages.success('Successfully deleted all issuers and keys');
      this.showDeleteAllIssuers = false;
      this.router.transitionTo('vault.cluster.secrets.backend.pki.configuration.index');
    } catch (error) {
      this.showDeleteAllIssuers = false;
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(message);
    }
  }
}
