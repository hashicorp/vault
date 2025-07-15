/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

import type FlagsService from 'vault/services/flags';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';

interface Args {
  onClose: () => void;
  onError: (errorMessage: string) => void;
  onConfirm: () => void;
}

export default class SyncActivationModal extends Component<Args> {
  @service declare readonly flags: FlagsService;
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;

  @tracked hasConfirmedDocs = false;

  @task
  @waitFor
  *onFeatureConfirm() {
    // clear any previous errors in the parent component
    this.args.onConfirm();

    // sync activation is managed by the root/administrative namespace so child namespaces are not sent.
    // for non-managed clusters the root namespace path is technically an empty string, otherwise we pass 'admin' if HVD managed.
    const namespace = this.flags.hvdManagedNamespaceRoot || '';
    try {
      yield this.api.sys.activationFlagsActivate_3(this.api.buildHeaders({ namespace }));
      // must refresh and not transition because transition does not refresh the model from within a namespace
      yield this.router.refresh('vault.cluster');
    } catch (error) {
      const { message } = yield this.api.parseError(error);
      this.args.onError(message);
      this.flashMessages.danger(`Error enabling feature \n ${message}`);
    } finally {
      this.args.onClose();
    }
  }
}
