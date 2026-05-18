/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import OidcAssignmentForm from 'vault/forms/oidc/assignment';
import {
  IdentityApiEntityListByIdListEnum,
  IdentityApiGroupListByIdListEnum,
} from '@hashicorp/vault-client-typescript';

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
  @service api;

  @tracked isFetching = true;
  @tracked entities = [];
  @tracked groups = [];

  constructor() {
    super(...arguments);
    this.form = new OidcAssignmentForm({ name: this.args.nameInput }, { isNew: true });
    this.fetchEntitiesAndGroups();
  }

  async fetchEntitiesAndGroups() {
    try {
      const [entitiesResult, groupsResult] = await Promise.allSettled([
        this.api.identity.entityListById(IdentityApiEntityListByIdListEnum.TRUE),
        this.api.identity.groupListById(IdentityApiGroupListByIdListEnum.TRUE),
      ]);

      this.entities =
        entitiesResult.status === 'fulfilled' ? this.api.keyInfoToArray(entitiesResult.value) : [];
      this.groups = groupsResult.status === 'fulfilled' ? this.api.keyInfoToArray(groupsResult.value) : [];
    } catch (error) {
      // swallow errors and render empty arrays
    }
    this.isFetching = false;
  }

  @action onSave(form) {
    this.args.onSave(form);
    this.form = null;
  }
}
