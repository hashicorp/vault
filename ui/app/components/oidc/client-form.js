import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import handleHasManySelection from 'core/utils/search-select-has-many';

/**
 * @module OidcClientForm
 * OidcClientForm components are used to create and update OIDC clients (a.k.a. applications)
 *
 * @example
 * ```js
 * <OidcClientForm @model={{this.model}} />
 * ```
 * @callback onCancel
 * @callback onSave
 * @param {Object} model - oidc client model
 * @param {onCancel} onCancel - callback triggered when cancel button is clicked
 * @param {onSave} onSave - callback triggered on save success
 */

export default class OidcClientForm extends Component {
  @service store;
  @service router;
  @service flashMessages;

  @tracked modelValidations;
  @tracked showMoreOptions = false;
  @tracked radioCardGroupValue = 'allow_all';

  @action
  async selectAssignments(selectedIds) {
    const assignments = await this.args.model.assignments;
    handleHasManySelection(selectedIds, assignments, this.store, 'oidc/assignment');
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isValid, state } = this.args.model.validate();
      this.modelValidations = isValid ? null : state;
      if (isValid) {
        yield this.args.model.save();
        this.flashMessages.success('Successfully created an application');
        this.args.onSave();
      }
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.flashMessages.danger(message);
    }
  }

  @action
  cancel() {
    const method = this.args.model.isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.args.model[method]();
    this.args.onCancel();
  }
}
