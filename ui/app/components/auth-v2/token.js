import Component from '@glimmer/component';
import { withAuthForm } from 'vault/decorators/auth-form';
import { task } from 'ember-concurrency';
import { next } from '@ember/runloop';
import errorMessage from 'vault/utils/error-message';

@withAuthForm('token')
export default class TokenComponent extends Component {
  constructor() {
    super(...arguments);
    if (this.args.wrappedToken) {
      next(() => {
        this.tryWrappedToken.perform();
      });
    }
  }

  @task *tryWrappedToken() {
    try {
      yield this.session.authenticate(
        `authenticator:token`,
        { token: this.args.wrappedToken },
        {
          backend: this.mountPath,
          namespace: this.args.namespace,
        }
      );
    } catch (e) {
      this.error = errorMessage(e);
    }
    if (this.session.isAuthenticated && this.args.onSuccess) {
      this.args.onSuccess();
    }
  }
}
