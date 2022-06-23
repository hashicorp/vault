import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
/**
 * @module OidcClientForm
 * OidcClientForm components are used to...
 *
 * @example
 * ```js
 * <OidcClientForm @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

export default class OidcClientForm extends Component {
  @service flashMessages;
  @service router;
  @tracked showMoreOptions = false;
  @tracked radioCardGroupValue = 'allow_all';

  @action
  async createClient() {
    try {
      this.args.model.save();
    } catch (e) {
      this.flashMessages.danger(e.errors?.join('. ') || e.message);
    }
  }

  @action
  cancel() {
    this.args.model.rollbackAttributes();
    this.router.transitionTo('vault.cluster.access.oidc.clients');
  }
}
