/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { capitalize } from '@ember/string';
import { task } from 'ember-concurrency';
import { addToArray } from 'vault/helpers/add-to-array';

export default class MfaMethodCreateController extends Controller {
  @service store;
  @service flashMessages;
  @service router;

  queryParams = ['type'];
  methods = [
    { name: 'TOTP', icon: 'history' },
    { name: 'Duo', icon: 'duo' },
    { name: 'Okta', icon: 'okta-color' },
    { name: 'PingID', icon: 'pingid' },
  ];

  @tracked type = null;
  @tracked method = null;
  @tracked enforcement;
  @tracked enforcementPreference = 'new';
  @tracked methodErrors;
  @tracked enforcementErrors;

  get description() {
    if (this.type === 'totp') {
      return `Once set up, TOTP requires a passcode to be presented alongside a Vault token when invoking an API request.
        The passcode will be validated against the TOTP key present in the identity of the caller in Vault.`;
    }
    return `Once set up, the ${this.formattedType} MFA method will require a push confirmation on mobile before login.`;
  }

  get formattedType() {
    if (!this.type) return '';
    return this.type === 'totp' ? this.type.toUpperCase() : capitalize(this.type);
  }
  get isTotp() {
    return this.type === 'totp';
  }
  get showForms() {
    return this.type && this.method;
  }

  @action
  onTypeSelect(type) {
    // set any form related properties to default values
    this.method = null;
    this.enforcement = null;
    this.methodErrors = null;
    this.enforcementErrors = null;
    this.enforcementPreference = 'new';
    this.type = type;
  }
  @action
  createModels() {
    if (this.method) {
      this.method.unloadRecord();
    }
    if (this.enforcement) {
      this.enforcement.unloadRecord();
    }
    this.method = this.store.createRecord('mfa-method', { type: this.type });
    this.enforcement = this.store.createRecord('mfa-login-enforcement');
  }
  @action
  onEnforcementPreferenceChange(preference) {
    if (preference === 'new') {
      this.enforcement = this.store.createRecord('mfa-login-enforcement');
    } else if (this.enforcement) {
      this.enforcement.unloadRecord();
      this.enforcement = null;
    }
    this.enforcementPreference = preference;
  }
  @action
  cancel() {
    this.method = null;
    this.enforcement = null;
    this.enforcementPreference = null;
    this.router.transitionTo('vault.cluster.access.mfa.methods');
  }
  @task
  *save() {
    const isValid = this.checkValidityState();
    if (isValid) {
      try {
        // first save method
        yield this.method.save();
        if (this.enforcement) {
          // mfa_methods is type PromiseManyArray so slice in necessary to convert it to an Array
          this.enforcement.mfa_methods = addToArray(this.enforcement.mfa_methods.slice(), this.method);
          try {
            // now save enforcement and catch error separately
            yield this.enforcement.save();
          } catch (error) {
            this.handleError(
              error,
              'Error saving enforcement. You can still create an enforcement separately and add this method to it.'
            );
          }
        }
        this.router.transitionTo('vault.cluster.access.mfa.methods.method', this.method.id);
      } catch (error) {
        this.handleError(error, 'Error saving method');
      }
    }
  }
  checkValidityState() {
    // block saving models if either is in an invalid state
    let isEnforcementValid = true;
    const methodValidations = this.method.validate();
    if (!methodValidations.isValid) {
      this.methodErrors = methodValidations.state;
    }
    // only validate enforcement if creating new
    if (this.enforcementPreference === 'new') {
      const enforcementValidations = this.enforcement.validate();
      // since we are adding the method after it has been saved ignore mfa_methods validation state
      const { name, targets } = enforcementValidations.state;
      isEnforcementValid = name.isValid && targets.isValid;
      if (!enforcementValidations.isValid) {
        this.enforcementErrors = enforcementValidations.state;
      }
    }
    return methodValidations.isValid && isEnforcementValid;
  }
  handleError(error, message) {
    const errorMessage = error?.errors ? `${message}: ${error.errors.join(', ')}` : message;
    this.flashMessages.danger(errorMessage);
  }
}
