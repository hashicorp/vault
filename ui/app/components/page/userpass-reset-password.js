import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import errorMessage from 'vault/utils/error-message';

export default class PageUserpassResetPasswordComponent extends Component {
  @service store;

  @tracked newPassword = '';
  @tracked successful = false;
  @tracked error = '';

  @action reset() {
    this.successful = false;
    this.error = '';
    this.newPassword = '';
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
      this.successful = true;
    } catch (e) {
      this.error = errorMessage(e, 'Check Vault logs for details');
    }
  }
}
