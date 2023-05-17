/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
// import errorMessage from 'vault/utils/error-message';

export default class PkiTidyStatusComponent extends Component {
  @service store;
  @service secretMountPath;
  @service flashMessages;

  @tracked tidyOptionsModal = false;
  @tracked confirmCancelTidy = false;
  @tracked tidyStatus = null;

  constructor(owner, args) {
    super(owner, args);

    const adapter = this.store.adapterFor('application');
    adapter
      .ajax(`/v1/${this.secretMountPath.currentPath}/tidy-status`, 'GET')
      .then(({ data }) => (this.tidyStatus = data));
  }

  get tidyStateAlertBanner() {
    let tidyState = this.tidyStatus?.state;

    if (this.cancelTidy.isRunning) {
      tidyState = 'Cancelling';
    } else if (this.cancelTidy.isSuccessful) {
      tidyState = 'Cancelled';
    }

    const tidyStateOptions = {
      Inactive: {
        color: 'highlight',
        title: 'Tidy is inactive',
        message: this.tidyStatus?.message,
      },
      Running: {
        color: 'highlight',
        title: 'Tidy in progress',
        message: this.tidyStatus?.message,
        shouldShowCancelTidy: true,
      },
      Finished: {
        color: 'success',
        title: 'Tidy operation finished',
        message: this.tidyStatus?.message,
      },
      Error: {
        color: 'warning',
        title: 'Tidy operation failed',
        message: this.tidyStatus?.message,
      },
      Cancelling: {
        color: 'warning',
        title: 'Tidy operation cancelling',
      },
      Cancelled: {
        color: 'warning',
        title: 'Tidy operation cancelled',
        message:
          'Your tidy operation has been cancelled. If this was a mistake configure and run another tidy operation.',
      },
    };

    return tidyStateOptions[tidyState];
  }

  @task
  @waitFor
  *cancelTidy() {
    try {
      const adapter = this.store.adapterFor('application');
      yield adapter.ajax(`/v1/${this.secretMountPath.currentPath}/tidy-cancel`, 'POST');
    } catch (e) {
      this.flashMessages.danger(e.errors.join(' '));
    }
  }
}
