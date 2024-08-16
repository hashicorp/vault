/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';

import type FlashMessageService from 'vault/services/flash-messages';
import type StoreService from 'vault/services/store';
import type RouterService from '@ember/routing/router-service';

interface Args {
  onClose: () => void;
  onError: (errorMessage: string) => void;
  onConfirm: () => void;
  isHvdManaged: boolean;
}

export default class SyncActivationModal extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly store: StoreService;
  @service declare readonly router: RouterService;

  @tracked hasConfirmedDocs = false;

  @task
  @waitFor
  *onFeatureConfirm() {
    // clear any previous errors in the parent component
    this.args.onConfirm();

    // must return null instead of root for non managed cluster.
    // child namespaces are not sent.
    const namespace = this.args.isHvdManaged ? 'admin' : null;

    try {
      yield this.store
        .adapterFor('application')
        .ajax('/v1/sys/activation-flags/secrets-sync/activate', 'POST', { namespace });
      // must refresh and not transition because transition does not refresh the model from within a namespace
      yield this.router.refresh('vault.cluster');
    } catch (error) {
      this.args.onError(errorMessage(error));
      this.flashMessages.danger(`Error enabling feature \n ${errorMessage(error)}`);
    } finally {
      this.args.onClose();
    }
  }
}
