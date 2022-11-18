import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

/**
 * @module ModalForm::OidcAssignmentTemplate
 * ModalForm::OidcAssignmentTemplate components are meant to render within a modal for creating a new assignment
 *
 * @example
 *  <ModalForm::OidcAssignmentTemplate
 *    @nameInput="new-item-name"
 *    @onSave={{this.closeModal}}
 *    @onCancel={{this.closeModal}}
 *  />
 * ```
 * @callback onCancel - callback triggered when cancel button is clicked
 * @callback onSave - callback triggered when save button is clicked
 * @param {string} nameInput - the name of the newly created policy
 */

export default class OidcAssignmentTemplate extends Component {
  @service store;
  @tracked assignment = null; // model record passed to oidc/assignment-form

  constructor() {
    super(...arguments);
    this.assignment = this.store.createRecord('oidc/assignment', { name: this.args.nameInput });
  }

  @action onSave(assignmentModel) {
    this.args.onSave(assignmentModel);
    // Reset component policy for next use
    this.assignment = null;
  }

  cleanup() {
    const method = this.assignment.isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.assignment[method]();
  }
}
