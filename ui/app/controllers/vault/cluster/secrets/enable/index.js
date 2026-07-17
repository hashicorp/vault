/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Controller from '@ember/controller';
import { action } from '@ember/object';

export default class SecretEnableController extends Controller {
  @service router;

  @action
  setMountType(type) {
    this.router.transitionTo('vault.cluster.secrets.enable.create', type);
  }
}
