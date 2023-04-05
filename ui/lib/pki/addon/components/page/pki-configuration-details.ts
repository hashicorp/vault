/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import errorMessage from 'vault/utils/error-message';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

//TYPES
import RouterService from '@ember/routing/router-service';
import FlashMessageService from 'vault/services/flash-messages';
import Store from '@ember-data/store';

interface Args {
  currentPath: string;
}

export default class PkiConfigurationDetails extends Component<Args> {
  @service declare readonly store: Store;
  @service declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked showDeleteAllIssuers = false;

  @action
  async deleteAllIssuers() {
    try {
      const issuerAdapter = this.store.adapterFor('pki/issuer');
      await issuerAdapter.deleteAllIssuers(this.args.currentPath);
      this.flashMessages.success('Successfully deleted all issuers and keys');
      this.showDeleteAllIssuers = false;
      this.router.transitionTo('vault.cluster.secrets.backend.pki.configuration.index');
      return;
    } catch (error) {
      this.showDeleteAllIssuers = false;
      this.flashMessages.danger(errorMessage(error));
    }
  }
}
