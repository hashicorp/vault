/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { capitalize } from '@ember/string';
import { task } from 'ember-concurrency';

import MfaCreateTotpMethodForm from 'vault/forms/mfa/method/totp';
import MfaCreateDuoMethodForm from 'vault/forms/mfa/method/duo';
import MfaCreateOktaMethodForm from 'vault/forms/mfa/method/okta';
import MfaCreatePingIdMethodForm from 'vault/forms/mfa/method/ping-id';
import MfaLoginEnforcementForm from 'vault/forms/mfa/login-enforcement';

export default class MfaMethodCreateController extends Controller {
  @service flashMessages;
  @service router;
  @service api;

  queryParams = ['type'];
  methods = [
    { name: 'TOTP', icon: 'history', type: 'totp' },
    { name: 'Duo', icon: 'duo-color', type: 'duo' },
    { name: 'Okta', icon: 'okta-color', type: 'okta' },
    { name: 'PingID', icon: 'ping-identity-color', type: 'pingid' },
  ];

  @tracked type = null;
  @tracked method = null;
  @tracked enforcement;
  @tracked enforcementPreference = 'new';
  @tracked methodErrors;
  @tracked enforcementErrors;

  get methodNameFromType() {
    return this.methods.find((method) => method.type === this.method.type)?.name || '';
  }

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
    return this.type === this.method?.type;
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
    if (this.type === 'totp') {
      this.method = new MfaCreateTotpMethodForm({ period: '30' }, { isNew: true });
    } else if (this.type === 'duo') {
      this.method = new MfaCreateDuoMethodForm({}, { isNew: true });
    } else if (this.type === 'okta') {
      this.method = new MfaCreateOktaMethodForm({}, { isNew: true });
    } else if (this.type === 'pingid') {
      this.method = new MfaCreatePingIdMethodForm({}, { isNew: true });
    }

    this.enforcement = new MfaLoginEnforcementForm({}, { isNew: true });
  }

  @action
  onEnforcementPreferenceChange(preference) {
    if (preference === 'new') {
      this.enforcement = new MfaLoginEnforcementForm({}, { isNew: true });
    } else if (this.enforcement) {
      this.enforcement = null;
    }
    this.enforcementPreference = preference;
  }
  @action
  onEnforcementSelect(enforcementForm) {
    this.enforcement = enforcementForm;
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
        const { data: methodData } = this.method.toJSON();

        let response = null;

        if (this.type === 'totp') {
          response = yield this.api.identity.mfaCreateTotpMethod({ ...methodData });
        } else if (this.type === 'duo') {
          response = yield this.api.identity.mfaCreateDuoMethod({ ...methodData });
        } else if (this.type === 'okta') {
          response = yield this.api.identity.mfaCreateOktaMethod({ ...methodData });
        } else if (this.type === 'pingid') {
          response = yield this.api.identity.mfaCreatePingIdMethod({ ...methodData });
        }

        const { data } = response;

        if (data.method_id && this.enforcement) {
          const { data: enforcementData } = this.enforcement.toJSON();
          // mfa_methods is type PromiseManyArray. Array methods like slice are no longer allowed on PromiseManyArray. We must yield the promise first, then call the method.
          // const mfaMethods = yield enforcementData.mfa_methods;
          // enforcementData.mfa_methods = addToArray(method_id, methodData);
          try {
            // now save enforcement and catch error separately
            yield this.api.identity.mfaWriteLoginEnforcement(enforcementData.name, {
              auth_method_accessors: enforcementData.auth_method_accessors || [],
              auth_method_types: enforcementData.auth_method_types || [],
              identity_entity_ids: enforcementData.identity_entity_ids || [],
              identity_group_ids: enforcementData.identity_group_ids || [],
              mfa_method_ids: [data.method_id],
            });
          } catch (error) {
            this.handleError(
              error,
              'Error saving enforcement. You can still create an enforcement separately and add this method to it.'
            );
          }
        }
        this.router.transitionTo('vault.cluster.access.mfa.methods.method', data.method_id);
      } catch (error) {
        this.handleError(error, 'Error saving method');
      }
    }
  }
  checkValidityState() {
    const { isValid: isMethodValid, state: methodState } = this.method.toJSON();
    // block saving models if either is in an invalid state
    let checkEnforcementValidity = true;

    if (!isMethodValid) {
      this.methodErrors = methodState;
    }
    // only validate enforcement if creating new and enforcement exists
    if (this.enforcementPreference === 'new' && this.enforcement) {
      const { isValid: isEnforcementValid, state: enforcementState } = this.enforcement.toJSON();
      const enforcementValidations = isEnforcementValid;
      // since we are adding the method after it has been saved ignore mfa_methods validation state
      const { name, targets } = enforcementState;
      checkEnforcementValidity = name.isValid && targets.isValid;
      if (!enforcementValidations.isValid) {
        this.enforcementErrors = enforcementState;
      }
    }
    return isMethodValid && checkEnforcementValidity;
  }
  handleError(error, message) {
    const errorMessage = error?.errors ? `${message}: ${error.errors.join(', ')}` : message;
    this.flashMessages.danger(errorMessage);
  }
}
