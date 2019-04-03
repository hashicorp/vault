import AuthConfigComponent from './config';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import DS from 'ember-data';

/**
 * @module AuthConfigForm/Options
 * The `AuthConfigForm/Options` is options portion of the auth config form.
 *
 * @example
 * ```js
 * {{auth-config-form/options model.model}}
 * ```
 *
 * @property model=null {String} - The corresponding auth model that is being configured.
 *
 */

export default AuthConfigComponent.extend({
  router: service(),
  wizard: service(),
  saveModel: task(function*() {
    let data = this.model.config.serialize();
    data.description = this.model.description;
    try {
      yield this.model.tune(data);
    } catch (err) {
      // AdapterErrors are handled by the error-message component
      // in the form
      if (err instanceof DS.AdapterError === false) {
        throw err;
      }
      return;
    }
    if (this.wizard.currentMachine === 'authentication' && this.wizard.featureState === 'config') {
      this.wizard.transitionFeatureMachine(this.wizard.featureState, 'CONTINUE');
    }
    this.router.transitionTo('vault.cluster.access.methods').followRedirects();
    this.flashMessages.success('The configuration was saved successfully.');
  }),
});
