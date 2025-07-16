/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

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
  @service store;
  @tracked key = null; // model record passed to oidc/key-form

  constructor() {
    super(...arguments);
    this.key = this.store.createRecord('oidc/key', { name: this.args.nameInput });
  }

  @action onSave(keyModel) {
    this.args.onSave(keyModel);
    // Reset component key for next use
    this.key = null;
  }
}
