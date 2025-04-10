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

import type FlagsService from 'vault/services/flags';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type Store from '@ember-data/store';

interface Args {
  onClose: () => void;
  onError: (errorMessage: string) => void;
  onConfirm: () => void;
}

export default class SyncActivationModal extends Component<Args> {
  @service declare readonly flags: FlagsService;
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly store: Store;

  @tracked hasConfirmedDocs = false;

  @task
  @waitFor
  *onFeatureConfirm() {
    // clear any previous errors in the parent component
    this.args.onConfirm();

    // sync activation is managed by the root/administrative namespace so child namespaces are not sent.
    // for non-managed clusters the root namespace path is technically an empty string so we pass null
    // otherwise we pass 'admin' if HVD managed.
    const namespace = this.flags.hvdManagedNamespaceRoot;

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
