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
 * @onClose={transition-to "vault.cluster.access.oidc.assignment"} @param1={{param1}}
 * @onSave={transition-to "vault.cluster.access.oidc.assignments.assignment.details" this.model.name}
 * />
 * ```
 * @param {object} model - The parent's model
 * @param {string} onClose - The transition-to helper directing where to go on cancel.
 * @param {string} onSave - The transition-to helper directing where to go on save.
 */

export default class OidcAssignmentForm extends Component {
  @service store;
  @service flashMessages;

  @tracked modelErrors;

  get errors() {
    return this.args.modelErrors || this.modelErrors;
  }

  @task
  *save() {
    this.modelErrors = {};
    // check validity state first and abort if invalid
    const { isValid, state } = this.args.model.validate();
    if (!isValid) {
      this.modelErrors = state;
    } else {
      try {
        yield this.args.model.save();
        this.args.onSave();
      } catch (error) {
        console.log(error, 'error');
        const message = error.errors ? error.errors.join('. ') : error.message;
        this.flashMessages.danger(message);
      }
    }
  }

  @action
  cancel() {
    // revert model changes
    this.args.model.rollbackAttributes();
    this.args.onClose();
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
    const groupIds = this.args.model.GroupIds;
    handleHasManySelection(selectedIds, groupIds, this.store, 'identity/group');
  }
}
