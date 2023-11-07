import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { withAuthForm } from 'vault/decorators/auth-form';

@withAuthForm('token')
export default class TokenComponent extends Component {
  @tracked token = '';

  @action
  handleToken(evt) {
    this.token = evt.target.value;
  }

  @action
  authenticate() {
    this.args.onAuthenticate({ token: this.token });
  }
}
