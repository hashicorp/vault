/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';

export default class MfaLoginEnforcementIndexController extends Controller {
  @service router;
  @service flashMessages;
  @service api;

  queryParams = ['tab'];
  tab = 'targets';

  @tracked showDeleteConfirmation = false;
  @tracked deleteError;

  @action
  async delete() {
    try {
      await this.api.identity.mfaDeleteLoginEnforcement(this.model.name);
      this.showDeleteConfirmation = false;
      this.flashMessages.success('MFA login enforcement deleted successfully');
      this.router.transitionTo('vault.cluster.access.mfa.enforcements');
    } catch (error) {
      this.deleteError = error;
    }
  }
}
