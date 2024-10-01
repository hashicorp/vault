/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

/**
 * @module ModalForm::OidcAssignmentTemplate
 * ModalForm::OidcAssignmentTemplate components render within a modal and create a model using the input from the search select. The model is passed to the oidc/assignment-form.
 *
 * @example
 *  <ModalForm::OidcAssignmentTemplate
 *    @nameInput="new-item-name"
 *    @onSave={{this.closeModal}}
 *    @onCancel={{@onCancel}}
 *  />
 * ```
 * @callback onCancel - callback triggered when cancel button is clicked
 * @callback onSave - callback triggered when save button is clicked
 * @param {string} nameInput - the name of the newly created assignment
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
    // Reset component assignment for next use
    this.assignment = null;
  }
}
