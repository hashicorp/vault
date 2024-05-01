/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

/**
 * @module ModalForm::PolicyTemplate
 * ModalForm::PolicyTemplate components are meant to render within a modal for creating a new policy of unknown type.
 *
 * @example
 *  <ModalForm::PolicyTemplate
 *    @nameInput="new-item-name"
 *    @onSave={{this.closeModal}}
 *    @onCancel={{this.closeModal}}
 *  />
 * ```
 * @callback onCancel - callback triggered when cancel button is clicked
 * @callback onSave - callback triggered when save button is clicked
 * @param {string} nameInput - the name of the newly created policy
 */

export default class PolicyTemplate extends Component {
  @service store;
  @service version;

  @tracked policy = null; // model record passed to policy-form

  get policyOptions() {
    return [
      { label: 'ACL Policy', value: 'acl', isDisabled: false },
      { label: 'Role Governing Policy', value: 'rgp', isDisabled: !this.version.hasSentinel },
    ];
  }

  @action
  setPolicyType(type) {
    if (this.policy) this.policy.unloadRecord(); // if user selects a different type, clear from store before creating a new record
    // Create form model once type is chosen
    this.policy = this.store.createRecord(`policy/${type}`, { name: this.args.nameInput });
  }

  @action
  onSave(policyModel) {
    this.args.onSave(policyModel);
    // Reset component policy for next use
    this.policy = null;
  }
}
