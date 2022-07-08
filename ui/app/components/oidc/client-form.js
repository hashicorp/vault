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
  @tracked radioCardGroupValue = 'allow_all';
  @tracked selectedAssignments = ['allow_all']; // default is 'allow_all' so no selected assignments

  constructor() {
    super(...arguments);
    this.fetchAssignments();
  }

  async fetchAssignments() {
    const assignmentIds = (await this.args.model.assignments).toArray().mapBy('id');
    this.radioCardGroupValue =
      assignmentIds.length === 0 || assignmentIds.includes('allow_all') ? 'allow_all' : 'limited';
  }

  @action handleSearchSelect(selectedIds) {
    this.selectedAssignments = this.radioCardGroupValue === 'allow_all' ? ['allow_all'] : selectedIds;
  }

  async selectModelAssignments() {
    const modelAssignments = await this.args.model.assignments;
    if (this.radioCardGroupValue === 'limited') {
      handleHasManySelection(
        this.selectedAssignments,
        modelAssignments,
        this.store,
        'oidc/assignment',
        this.args.model
      );
    } else {
      // search select hasn't queried the assignments unless the "limit" radio select has been interacted with
      // so need to make a network request to fetch allow_all record
      // move to init?
      const allowAllRecord = await this.store.findRecord('oidc/assignment', 'allow_all');
      modelAssignments.addObject(allowAllRecord);
      this.args.model.save();
    }
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isValid, state } = this.args.model.validate();
      this.modelValidations = isValid ? null : state;
      if (isValid) {
        this.selectModelAssignments();
        yield this.args.model.save();
        this.flashMessages.success(
          `Successfully ${this.args.model.isNew ? 'created an' : 'updated'} application`
        );
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
