import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import handleHasManySelection from 'core/utils/search-select-has-many';

/**
 * @module Oidc::AssignmentForm
 * Oidc::AssignmentForm components are used to display the create view for OIDC providers assignments.
 *
 * @example
 * ```js
 * <Oidc::AssignmentForm @model={this.model}
 * @onCancel={transition-to "vault.cluster.access.oidc.assignment"} @param1={{param1}}
 * @onSave={transition-to "vault.cluster.access.oidc.assignments.assignment.details" this.model.name}
 * />
 * ```
 * @callback onCancel
 * @callback onSave
 * @param {object} model - The parent's model
 * @param {string} onCancel - callback triggered when cancel button is clicked
 * @param {string} onSave - callback triggered when save button is clicked
 */

export default class OidcAssignmentFormComponent extends Component {
  @service store;
  @service flashMessages;

  @tracked modelValidations;

  @task
  *save() {
    event.preventDefault();
    try {
      const { isValid, state } = this.args.model.validate();
      this.modelValidations = isValid ? null : state;
      if (isValid) {
        yield this.args.model.save();
        this.flashMessages.success('Successfully created an assignment');
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

  @action
  handleOperation(e) {
    let value = e.target.value;
    this.args.model.name = value;
  }

  @action
  onEntitiesSelect(selectedIds) {
    const entityIds = this.args.model.entityIds;
    handleHasManySelection(selectedIds, entityIds, this.store, 'identity/entity');
  }

  @action
  onGroupsSelect(selectedIds) {
    const groupIds = this.args.model.groupIds;
    handleHasManySelection(selectedIds, groupIds, this.store, 'identity/group');
  }
}
