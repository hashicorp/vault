/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import errorMessage from 'vault/utils/error-message';

export default class PageUserpassResetPasswordComponent extends Component {
  @service store;
  @service flashMessages;

  @tracked newPassword = '';
  @tracked error = '';

  onSuccess() {
    this.error = '';
    this.newPassword = '';
    this.flashMessages.success('Successfully reset password');
  }

  @task
  *updatePassword(evt) {
    evt.preventDefault();
    this.error = '';
    const adapter = this.store.adapterFor('auth-method');
    const { backend, username } = this.args;
    if (!backend || !username) return;
    if (!this.newPassword) {
      this.error = 'Please provide a new password.';
      return;
    }
    try {
      yield adapter.resetPassword(backend, username, this.newPassword);
      this.onSuccess();
    } catch (e) {
      this.error = errorMessage(e, 'Check Vault logs for details');
    }
  }
}
