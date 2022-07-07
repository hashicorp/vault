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
  @service flashMessages;

  @tracked modelValidations;
  @tracked radioCardGroupValue;
  @tracked selectedAssignments;

  constructor() {
    super(...arguments);
    this.fetchAssignments();
  }

  async fetchAssignments() {
    const assignments = (await this.args.model.assignments).toArray().mapBy('id');
    this.radioCardGroupValue =
      assignments.length === 0 || assignments.includes('allow_all') ? 'allow_all' : 'limited';
  }

  @action
  onChange(selection) {
    if (typeof selection === 'string') this.radioCardGroupValue = selection;
    if (Array.isArray(selection)) {
      this.selectedAssignments = selection;
    }
  }

  @action
  async handleAssignmentSelection() {
    const assignments = await this.args.model.assignments;
    console.log(assignments, 'assignments');
    let selection = this.radioCardGroupValue === 'allow_all' ? ['allow_all'] : this.selectedAssignments;
    handleHasManySelection(selection, assignments.toArray(), this.store, 'oidc/assignment');
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      this.handleAssignmentSelection();
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
