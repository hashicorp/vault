/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Controller from '@ember/controller';
import { service } from '@ember/service';

export default class SealController extends Controller {
  @service api;
  @service auth;
  @service router;
  @service version;

  @action
  seal() {
    return this.api.request.put('/sys/seal').then(() => {
      this.router.transitionTo('vault.cluster.unseal');
      this.version.fetchType();
      this.auth.deleteCurrentToken();
      // Reset version so it doesn't show on footer
      this.version.version = null;
    });
  }
}
