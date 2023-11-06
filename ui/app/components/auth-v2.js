import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';

export default class AuthV2Component extends Component {
  @service session;
  @tracked token = '';
  @tracked error = '';
  type = 'token';

  @action
  handleChange(evt) {
    this.token = evt.target.value;
  }

  @action
  async handleLogin(evt) {
    evt.preventDefault();
    const authenticator = `authenticator:${this.type}`;
    try {
      await this.session.authenticate(authenticator, this.token);
    } catch (e) {
      this.error = errorMessage(e);
    }

    if (this.session.isAuthenticated) {
      this.args.onSuccess();
    }
  }
}
