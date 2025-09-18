/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class VaultClusterMfaSetupController extends Controller {
  @service auth;
  @tracked onStep = 1;
  @tracked warning = '';
  @tracked uuid = '';
  @tracked qrCode = '';

  header = 'MFA Setup';
  description =
    'TOTP Multi-factor authentication (MFA) can be enabled here if it is required by your administrator. This will ensure that you are not prevented from logging into Vault in the future, once MFA is fully enforced.';

  get entityId() {
    return this.auth.authData.entityId;
  }

  @action isUUIDVerified(verified) {
    this.warning = ''; // clear the warning, otherwise it persists.
    if (verified) {
      this.onStep = 2;
    } else {
      this.restartFlow();
    }
  }

  @action
  restartFlow() {
    this.onStep = 1;
  }

  @action
  saveUUIDandQrCode(uuid, qrCode) {
    // qrCode could be an empty string if the admin-generate was not successful
    this.uuid = uuid;
    this.qrCode = qrCode;
  }

  @action
  showWarning(warning) {
    this.warning = warning;
    this.onStep = 2;
  }
}
