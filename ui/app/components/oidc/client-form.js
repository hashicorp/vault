import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
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
  @tracked selectedAssignments;
  @tracked radioCardGroupValue =
    !this.args.model.assignments || this.args.model.assignments.includes('allow_all')
      ? 'allow_all'
      : 'limited';

  @action
  handleAssignmentSelection(selection) {
    this.selectedAssignments = this.radioCardGroupValue === 'allow_all' ? [] : selection;
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
          // to avoid having 'allow_all' appear in the "limited" search-select dropdown
          // the UI limits the config by allowing either 'allow_all' OR a list of other assignments
          // note: when editing the UI removes any additional assignments previously configured via CLI
          this.args.model.assignments = ['allow_all'];
        } else {
          this.args.model.assignments = this.selectedAssignments;
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
