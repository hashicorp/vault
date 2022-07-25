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
  @tracked selectedAssignments;

  constructor() {
    super(...arguments);
    this.fetchAndSetAssignments();
  }

  async fetchAndSetAssignments() {
    const assignmentIds = (await this.args.model.assignments).toArray().mapBy('id');
    this.selectedAssignments = assignmentIds;
    this.radioCardGroupValue =
      assignmentIds.length === 0 || assignmentIds.includes('allow_all') ? 'allow_all' : 'limited';
  }

  @action
  // handle radio buttons and search-select changes here
  async handleAssignmentSelection(selection) {
    const modelAssignments = await this.args.model.assignments;
    if (typeof selection === 'string') {
      // handle radio buttons
      this.radioCardGroupValue = selection;
      // clear out selected assignments, do not set it to ['allow_all']
      // because it will appear in the search select if the user toggles back to "limited"
      this.selectedAssignments = [];
    }
    if (Array.isArray(selection)) {
      // handle search select
      this.selectedAssignments = selection;
    }
    handleHasManySelection(this.selectedAssignments, modelAssignments, this.store, 'oidc/assignment');
  }

  @action
  async allowAllAssignments() {
    const modelAssignments = await this.args.model.assignments;
    const allowAllModel = await this.store.findRecord('oidc/assignment', 'allow_all');
    modelAssignments.addObject(allowAllModel);
    this.selectedAssignments = ['allow_all'];
    handleHasManySelection(this.selectedAssignments, modelAssignments, this.store, 'oidc/assignment');
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isValid, state } = this.args.model.validate();
      this.modelValidations = isValid ? null : state;
      if (isValid) {
        if (this.radioCardGroupValue === 'allow_all') {
          // the backend permits 'allow_all' AND other assignments, though 'allow_all' will take precedence
          // this is slightly awkward UX when 'allow_all' appears in the "limited" search-select dropdown
          // so the UI limits the config by allowing either 'allow_all' OR a list of other assignments
          yield this.allowAllAssignments();
        }
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
