import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import handleHasManySelection from 'core/utils/search-select-has-many';

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
  @service store;
  @service router;
  @service flashMessages;
  @tracked showMoreOptions = false;
  @tracked radioCardGroupValue = 'allow_all';

  @action
  async selectAssignments(selectedIds) {
    const assignments = await this.args.model.assignments;
    handleHasManySelection(selectedIds, assignments, this.store, 'oidc/assignment');
  }

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
