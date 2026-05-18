/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import PolicyForm from 'vault/forms/policy';

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
  @tracked form = null; // form class passed to policy-form

  get policyOptions() {
    return [
      { label: 'ACL Policy', value: 'acl', isDisabled: false },
      { label: 'Role Governing Policy', value: 'rgp', isDisabled: !this.version.hasSentinel },
    ];
  }

  @action
  setPolicyType(type) {
    // Create form once type is chosen
    const policyForm = new PolicyForm(
      { name: this.args.nameInput, enforcement_level: 'hard-mandatory' },
      { isNew: true }
    );
    policyForm.policyType = type;
    this.form = policyForm;
  }
  @action
  onSave(policyForm) {
    this.args.onSave(policyForm);
    // Reset component policy for next use
    this.policy = null;
  }
}
