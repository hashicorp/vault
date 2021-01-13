import Component from '@ember/component';
import { inject as service } from '@ember/service';
/**
 * @module GetCredentialsCard
//  * ARG TODO update text here
 * SelectableCard components are card-like components that display a title, total, subtotal, and anything after the yield.
 * They are designed to be used in containers that act as flexbox or css grid containers.
 *
 * @example
 * ```js
 * <GetCredentialsCard @title="Get Credentials" @models={{array 'database/roles'}} />
 * ```
 * @param title=null {String} - cardTitle displays the card title
 * @param models=null {Array} - An array of model types to fetch from the API.  Passed through to SearchSelect component
 */

// ARG TODO turn into octane component
// ARG TODO add in remaining params and storybook this.
export default Component.extend({
  tagName: '', // do not wrap component with div
  router: service(),
  role: '',
  buttonDisabled: true,
  actions: {
    getSelectedValue(selectValue) {
      this.role = selectValue[0];
      this.toggleProperty('buttonDisabled');
    },
    transitionToCredential() {
      let role = this.role;
      this.router.transitionTo('vault.cluster.secrets.backend.credentials', role);
    },
  },
});
