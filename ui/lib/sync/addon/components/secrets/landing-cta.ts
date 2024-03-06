/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';

import type VersionService from 'vault/services/version';
import type RouterService from '@ember/routing/router-service';
import type StoreService from 'vault/services/store';
import type FlashMessageService from 'vault/services/flash-messages';
import type { HTMLElementEvent } from 'vault/forms';

interface Args {
  featureEnabled: boolean;
}

export default class SyncLandingCtaComponent extends Component<Args> {
  @service declare readonly version: VersionService;
  @service declare readonly router: RouterService;
  @service declare readonly store: StoreService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked showEnableFeatureModal = false;
  @tracked confirmDisabled = true;

  createDestination() {
    this.router.transitionTo('vault.cluster.sync.secrets.destinations.create');
  }

  @action
  onCreate() {
    if (this.args.featureEnabled) {
      this.createDestination();
    } else {
      this.showEnableFeatureModal = true;
    }
  }

  @action
  onDocsConfirmChange(event: HTMLElementEvent<HTMLInputElement>) {
    this.confirmDisabled = !event.target.checked;
  }

  @task
  @waitFor
  *onFeatureConfirm() {
    try {
      // payload is empty
      const payload = { data: {} };
      yield this.store.adapterFor('application').ajax('/v1/sys/activation-flags/secrets-sync/activate', 'PATCH', payload);
      this.createDestination();
    } catch (error) {
      this.flashMessages.danger(`Error enabling feature \n ${errorMessage(error)}`);
    }
  }
}
