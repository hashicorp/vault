/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';

export default class PageUserpassResetPasswordComponent extends Component {
  @service flashMessages;
  @service api;

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
    const { backend, username } = this.args;
    if (!backend || !username) return;
    if (!this.newPassword) {
      this.error = 'Please provide a new password.';
      return;
    }
    try {
      yield this.api.auth.userpassResetPassword(username, backend, { password: this.newPassword });
      this.onSuccess();
    } catch (e) {
      const { message } = yield this.api.parseError(e, 'Check Vault logs for details');
      this.error = message;
    }
  }
}
