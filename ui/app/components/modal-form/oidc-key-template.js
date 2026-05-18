/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import OidcKeyForm from 'vault/forms/oidc/key';

/**
 * @module ModalForm::OidcKeyTemplate
 * ModalForm::OidcKeyTemplate components render within a modal and create a model using the input from the search select. The model is passed to the oidc/key-form.
 *
 * @example
 *  <ModalForm::OidcKeyTemplate
 *    @nameInput="new-key-name"
 *    @onSave={{this.closeModal}}
 *    @onCancel={{@onCancel}}
 *  />
 *
 * @callback onCancel - callback triggered when cancel button is clicked
 * @callback onSave - callback triggered when save button is clicked
 * @param {string} nameInput - the name of the newly created key
 */

export default class OidcKeyTemplate extends Component {
  constructor() {
    super(...arguments);
    this.form = new OidcKeyForm({ name: this.args.nameInput }, { isNew: true });
  }

  @action onSave(form) {
    this.args.onSave(form);
    this.form = null;
  }
}
