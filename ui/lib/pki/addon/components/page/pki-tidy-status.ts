/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import type Store from '@ember-data/store';
import type PkiTidyModel from 'vault/models/pki/tidy';
import type SecretMountPath from 'vault/services/secret-mount-path';

interface Args {
  autoTidyConfig: PkiTidyModel;
}
interface TidyStatus {
  data: {
    safety_buffer: number;
    tidy_cert_store: boolean;
    tidy_revoked_certs: boolean;
    state: string;
    error: string;
    time_started: string | null;
    time_finished: string | null;
    message: string;
    cert_store_deleted_count: number;
    revoked_cert_deleted_count: number;
    missing_issuer_cert_count: number;
    tidy_expired_issuers: boolean;
    issuer_safety_buffer: string;
    tidy_move_legacy_ca_bundle: boolean;
    tidy_revocation_queue: boolean;
    revocation_queue_deleted_count: number;
    tidy_cross_cluster_revoked_certs: boolean;
    cross_revoked_cert_deleted_count: number;
  };
}
export default class PkiTidyStatusComponent extends Component<Args> {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  @tracked tidyOptionsModal = false;
  @tracked confirmCancelTidy = false;
  @tracked tidyStatus = {};

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    const adapter = this.store.adapterFor('application');
    adapter
      .ajax(`/v1/${this.secretMountPath.currentPath}/tidy-status`, 'GET')
      .then((res: TidyStatus) => (this.tidyStatus = res.data));
  }

  // get tidyState() {

  // }

  @action
  cancelTidy() {
    // do the thing
  }
}
