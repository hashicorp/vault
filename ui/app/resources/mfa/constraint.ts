/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { ParsedMfaConstraint, ParsedMfaMethod } from 'vault/vault/auth/mfa';

/* 
There is a slight difference between a login enforcement configuration and the data returned in the `mfa_constraints` key when a user is logging in.
Although there is overlap, there are slight differences between the two and are referred to as MfaConstraint or MfaLoginEnforcement to differentiate between the two.
 * "MfaConstraint" refers to the enforcement data in the `mfa_constraints` key
    - the `self_enrollment_enabled` boolean for example is flipped by the backend if a user has already enrolled, 
      meaning it doesn't just return the config value but whether or not the UI should enter the self-enroll workflow.
 * "MfaLoginEnforcement" is reserved for parameters used when configuring the enforcement. (see https://developer.hashicorp.com/vault/api-docs/secret/identity/mfa/login-enforcement#create-a-login-enforcement)
  (The MfaLoginEnforcement resource will exist when we migrate away from the ember data model mfa-login-enforcement )
*/

export default class MfaConstraint {
  @tracked name: string;
  @tracked methods: ParsedMfaMethod[] = [];
  @tracked selectedMethod: ParsedMfaMethod | undefined = undefined;
  // These are set on the login enforcement and not the MFA method
  // because they only correspond to the selectedMethod.
  @tracked passcode = '';
  @tracked qrCode = '';

  constructor(constraint: ParsedMfaConstraint) {
    const { name, methods } = constraint;
    this.name = name;
    this.methods = methods;
    // Only set selectedMethod if there is only one method.
    // Otherwise, the user must select which method they want to verify with.
    this.selectedMethod = this.methods.length === 1 ? methods[0] : undefined;
  }

  @action
  setPasscode(value: string) {
    this.passcode = value;
  }

  @action
  setSelectedMethod(id: string) {
    const method = this.methods.find((m) => m.id === id);
    this.selectedMethod = method;
  }

  get hasSelfEnrollMethods(): boolean {
    return !!this.selfEnrollMethods.length;
  }

  // To be "true", the login enforcement must have a selected method set
  // and passcode, if applicable.
  get isSatisfied() {
    if (this.selectedMethod) {
      return this.selectedMethod.uses_passcode ? !!this.passcode : true;
    }
    return false;
  }

  get selfEnrollMethods(): ParsedMfaMethod[] | [] {
    // Self-enrollment is an enterprise only feature and self_enrollment_enabled will always be false for CE
    // It also returns false if the user already has an MFA secret (meaning they have already enrolled.)
    return this.methods.filter((m) => m?.self_enrollment_enabled);
  }

  get validateData() {
    return {
      methods: this.methods,
      passcode: this.passcode,
      selectedMethod: this.selectedMethod,
    };
  }
}
