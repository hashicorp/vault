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

import type Store from '@ember-data/store';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type FlashMessageService from 'vault/services/flash-messages';
import type VersionService from 'vault/services/version';
import type PkiTidyModel from 'vault/models/pki/tidy';
import type RouterService from '@ember/routing/router-service';

interface Args {
  autoTidyConfig: PkiTidyModel;
  tidyStatus: TidyStatusParams;
}

interface TidyStatusParams {
  // tidy banner
  state: string;
  error: string;
  message: string;
  // tidy status
  time_started: string | null;
  time_finished: string | null;
  cert_store_deleted_count: number;
  revoked_cert_deleted_count: number;
  missing_issuer_cert_count: number;
  revocation_queue_deleted_count: number; // enterprise only
  cross_revoked_cert_deleted_count: number; // enterprise only
  // tidy settings
  tidy_cert_store: boolean;
  tidy_revoked_certs: boolean;
  tidy_expired_issuers: boolean;
  safety_buffer: number;
  tidy_move_legacy_ca_bundle: boolean;
  issuer_safety_buffer: string;
  tidy_revocation_queue: boolean; // enterprise only
  tidy_cross_cluster_revoked_certs: boolean; // enterprise only
  revocation_queue_safety_buffer: string; // enterprise only
}

export default class PkiTidyStatusComponent extends Component<Args> {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly version: VersionService;
  @service declare readonly router: RouterService;

  @tracked tidyOptionsModal = false;
  @tracked confirmCancelTidy = false;

  tidyStatusGeneralFields = [
    'time_started',
    'time_finished',
    'cert_store_deleted_count',
    'revoked_cert_deleted_count',
    'missing_issuer_cert_count',
  ];

  tidyStatusConfigFields = [
    'tidy_cert_store',
    'tidy_cert_metadata',
    'tidy_revoked_certs',
    'safety_buffer',
    'pause_duration',
    'tidy_expired_issuers',
    'tidy_move_legacy_ca_bundle',
    'issuer_safety_buffer',
  ];

  // enterprise only
  crossClusterOperation = {
    status: ['revocation_queue_deleted_count', 'cross_revoked_cert_deleted_count'],
    config: ['tidy_revocation_queue', 'tidy_cross_cluster_revoked_certs', 'revocation_queue_safety_buffer'],
  };

  get isEnterprise() {
    return this.version.isEnterprise;
  }

  get tidyState() {
    return this.args.tidyStatus?.state;
  }

  get hasTidyConfig() {
    return !this.tidyStatusConfigFields.every(
      (attr) => this.args.tidyStatus[attr as keyof TidyStatusParams] === null
    );
  }

  get tidyStateAlertBanner() {
    const tidyState = this.tidyState;

    switch (tidyState) {
      case 'Inactive':
        return {
          color: 'highlight',
          title: 'Tidy is inactive',
          message: this.args.tidyStatus?.message,
        };
      case 'Running':
        return {
          color: 'highlight',
          title: 'Tidy in progress',
          message: this.args.tidyStatus?.message,
          shouldShowCancelTidy: true,
        };
      case 'Finished':
        return {
          color: 'success',
          title: 'Tidy operation finished',
          message: this.args.tidyStatus?.message,
        };
      case 'Error':
        return {
          color: 'warning',
          title: 'Tidy operation failed',
          message: this.args.tidyStatus?.error,
        };
      case 'Cancelling':
        return {
          color: 'warning',
          title: 'Tidy operation cancelling',
          icon: 'loading',
        };
      case 'Cancelled':
        return {
          color: 'warning',
          title: 'Tidy operation cancelled',
          message:
            'Your tidy operation has been cancelled. If this was a mistake configure and run another tidy operation.',
        };
      default:
        return {
          color: 'warning',
          title: 'Tidy status not found',
          message: "The system reported no tidy status. It's recommended to perform a new tidy operation.",
        };
    }
  }

  @task
  @waitFor
  *cancelTidy() {
    try {
      const tidyAdapter = this.store.adapterFor('pki/tidy');
      yield tidyAdapter.cancelTidy(this.secretMountPath.currentPath);
      this.router.transitionTo('vault.cluster.secrets.backend.pki.tidy');
    } catch (error) {
      this.flashMessages.danger(errorMessage(error));
    } finally {
      this.confirmCancelTidy = false;
    }
  }
}
