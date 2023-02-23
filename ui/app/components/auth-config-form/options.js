import AdapterError from '@ember-data/adapter/error';
import AuthConfigComponent from './config';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

/**
 * @module AuthConfigForm/Options
 * The `AuthConfigForm/Options` is options portion of the auth config form.
 *
 * @example
 * ```js
 * {{auth-config-form/options model.model}}
 * ```
 *
 * @property model=null {DS.Model} - The corresponding auth model that is being configured.
 *
 */

export default AuthConfigComponent.extend({
  flashMessages: service(),
  router: service(),

  saveModel: task(
    waitFor(function* () {
      const data = this.model.config.serialize();
      data.description = this.model.description;

      // token_type should not be tuneable for the token auth method.
      if (this.model.type === 'token') {
        delete data.token_type;
      }

      try {
        yield this.model.tune(data);
      } catch (err) {
        // AdapterErrors are handled by the error-message component
        // in the form
        if (err instanceof AdapterError === false) {
          throw err;
        }
        // because we're not calling model.save the model never updates with
        // the error.  Forcing the error message by manually setting the errorMessage
        try {
          this.model.set('errorMessage', err.errors.firstObject);
        } catch {
          // do nothing
        }
        return;
      }
      this.router.transitionTo('vault.cluster.access.methods').followRedirects();
      this.flashMessages.success('The configuration was saved successfully.');
    })
  ),
});
