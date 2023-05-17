/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

export default class PkiTidyStatusComponent extends Component {
  @service store;
  @service secretMountPath;
  @service flashMessages;
  @service router;

  @tracked tidyOptionsModal = false;
  @tracked confirmCancelTidy = false;

  get generalSectionFields() {
    return [
      'time_started',
      'time_finished',
      'last_auto_tidy_finished',
      'cert_store_deleted_count',
      'missing_issuer_cert_count',
      'revocation_queue_deleted_count',
    ];
  }
  get universalSectionFields() {
    return [
      'tidy_cert_store',
      'tidy_revocation_queue',
      'tidy_cross_cluster_revoked_certs',
      'safety_buffer',
      'pause_duration',
    ];
  }
  get issuersSectionFields() {
    return ['tidy_expired_issuers', 'tidy_move_legacy_ca_bundle', 'issuer_safety_buffer'];
  }
  get crossClusterOperation() {
    return ['tidy_revocation_queue', 'revocation_queue_safety_buffer'];
  }

  get tidyStateAlertBanner() {
    let tidyState = this.args.tidyStatus?.state;

    if (this.cancelTidy.isRunning) {
      tidyState = 'Cancelling';
    } else if (this.cancelTidy.isSuccessful) {
      tidyState = 'Cancelled';
    }

    const tidyStateOptions = {
      Inactive: {
        color: 'highlight',
        title: 'Tidy is inactive',
        message: this.args.tidyStatus?.message,
      },
      Running: {
        color: 'highlight',
        title: 'Tidy in progress',
        message: this.args.tidyStatus?.message,
        shouldShowCancelTidy: true,
      },
      Finished: {
        color: 'success',
        title: 'Tidy operation finished',
        message: this.args.tidyStatus?.message,
      },
      Error: {
        color: 'warning',
        title: 'Tidy operation failed',
        message: this.args.tidyStatus?.message,
      },
      Cancelling: {
        color: 'warning',
        title: 'Tidy operation cancelling',
        icon: 'loading',
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
    // TODO: make a custom adapter method when Claire merges her form work!
    // TODO: fix the transition from cancelling to cancelled state.
    try {
      const adapter = this.store.adapterFor('application');
      yield adapter.ajax(`/v1/${this.secretMountPath.currentPath}/tidy-cancel`, 'POST');
      this.router.transitionTo('vault.cluster.secrets.backend.pki.tidy.index');
    } catch (e) {
      this.flashMessages.danger(e.errors.join(' '));
    }
  }
}
