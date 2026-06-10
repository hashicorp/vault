/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { addToArray } from 'vault/helpers/add-to-array';
import { removeFromArray } from 'vault/helpers/remove-from-array';
import AuthMethodResource from 'vault/resources/auth/method';
import { prepareTargets } from 'vault/utils/mfa-login-enforcement-helpers';

/**
 * @module MfaLoginEnforcementForm
 * MfaLoginEnforcementForm components are used to create and edit login enforcements
 *
 * @example
 * ```js
 * <MfaLoginEnforcementForm @model={{this.model}} @isInline={{false}} @onSave={{this.onSave}} @onClose={{this.onClose}} />
 * ```
 * @callback onSave
 * @callback onClose
 * @param {Object} model - login enforcement form model
 * @param {Object} [isInline] - toggles inline display of form -- method selector and actions are hidden and should be handled externally
 * @param {Object} [modelErrors] - model validations state object if handling actions externally when displaying inline
 * @param {onSave} [onSave] - triggered on save success
 * @param {onClose} [onClose] - triggered on cancel
 */

export default class MfaLoginEnforcementForm extends Component {
  @service flashMessages;
  @service api;

  targetTypes = [
    { label: 'Authentication mount', type: 'accessor', key: 'auth_method_accessors' },
    { label: 'Authentication method', type: 'method', key: 'auth_method_types' },
    { label: 'Group', type: 'identity/group', key: 'identity_group_ids' },
    { label: 'Entity', type: 'identity/entity', key: 'identity_entity_ids' },
  ];

  searchSelectOptions = null;

  @tracked name;
  @tracked targets = [];
  @tracked selectedTargetType = 'accessor';
  @tracked selectedTargetValue = null;
  @tracked searchSelect = {
    options: [],
    selected: [],
  };
  @tracked authMethods = [];
  @tracked modelErrors;

  mfaMethods = []; // does not change after initial fetch, thus not tracked

  constructor() {
    super(...arguments);
    // aggregate different target array properties on model into flat list
    this.flattenTargets();
    // eagerly fetch identity groups and entities for use as search select options
    this.resetTargetState();
    // only auth method types that have mounts can be selected as targets -- fetch from sys/auth and map by type
    this.fetchAuthMethods();
    this.fetchMfaMethods();
  }

  async flattenTargets() {
    const preparedTargets = await prepareTargets(this.args.form, this.api, {
      includeFormFields: true,
    });
    this.targets = preparedTargets;
  }

  async resetTargetState() {
    this.selectedTargetValue = null;
    const options = this.searchSelectOptions || {};
    if (!this.searchSelectOptions) {
      const types = ['identity/group', 'identity/entity'];
      for (const type of types) {
        try {
          const apiMethod = type === 'identity/group' ? 'groupListById' : 'entityListById';
          const response = await this.api.identity[apiMethod](true);
          options[type] = this.api.keyInfoToArray(response);
        } catch (error) {
          options[type] = [];
        }
      }
      this.searchSelectOptions = options;
    }
    if (this.selectedTargetType.includes('identity')) {
      this.searchSelect = {
        selected: [],
        options: [...options[this.selectedTargetType]],
      };
    }
  }
  async fetchAuthMethods() {
    const { data } = await this.api.sys.authListEnabledMethods();
    this.authMethods = this.api
      .responseObjectToArray(data, 'path')
      .map((method) => new AuthMethodResource(method, this).methodType)
      .uniq();
  }

  async fetchMfaMethods() {
    // @methods now contains API response data (plain objects) instead of Ember Data models
    const methods = this.args.methods || [];

    // If the form already has mfa_methods set (e.g., editing existing enforcement),
    // filter to show only the selected methods
    if (this.args.form.mfa_methods && this.args.form.mfa_methods.length > 0) {
      this.mfaMethods = methods.filter((method) => this.args.form.mfa_methods.includes(method.id));
    } else {
      // For new enforcements, start with empty selection
      this.mfaMethods = [];
    }
  }

  get selectedTarget() {
    return this.targetTypes.find((tt) => tt.type === this.selectedTargetType);
  }

  get errors() {
    return this.args.modelErrors || this.modelErrors;
  }

  updateModelForKey(key) {
    // For identity entities and groups, extract IDs from the model objects
    if (key === 'identity_entity_ids' || key === 'identity_group_ids') {
      const newValue = this.targets.filter((t) => t.key === key).map((t) => t.value.id);
      this.args.form[key] = newValue;
    } else {
      const newValue = this.targets.filter((t) => t.key === key).map((t) => t.value);
      this.args.form[key] = newValue;
    }
  }

  @task
  *save() {
    this.modelErrors = {};
    // check validity state first and abort if invalid
    const { data, isValid, state } = this.args.form.toJSON();

    if (!isValid) {
      this.modelErrors = state;
    } else {
      try {
        const { name, mfa_methods, ...enforcementData } = data;
        yield this.api.identity.mfaWriteLoginEnforcement(name, {
          ...enforcementData,
          auth_method_accessors: enforcementData.auth_method_accessors || [],
          auth_method_types: enforcementData.auth_method_types || [],
          identity_entity_ids: enforcementData.identity_entity_ids || [],
          identity_group_ids: enforcementData.identity_group_ids || [],
          mfa_method_ids: mfa_methods || [],
        });
        this.args.onSave();
      } catch (error) {
        const message = error.errors ? error.errors.join('. ') : error.message;
        this.flashMessages.danger(message);
      }
    }
  }

  @action
  async onMethodChange(selectedIds) {
    // Update the form's mfa_methods field with the selected IDs
    this.args.form.mfa_methods = selectedIds;

    // Update mfaMethods to reflect the current selection for the SearchSelect component
    const methods = this.args.methods || [];
    this.mfaMethods = methods.filter((method) => selectedIds.includes(method.id));
  }

  @action
  onTargetSelect(type) {
    this.selectedTargetType = type;
    this.resetTargetState();
  }
  @action
  setTargetValue(selected) {
    const { type } = this.selectedTarget;
    if (type.includes('identity')) {
      // Find the selected item from the already-fetched options
      const selectedItem = this.searchSelectOptions[type].find((item) => item.id === selected[0]);
      this.selectedTargetValue = selectedItem;
    } else {
      this.selectedTargetValue = selected;
    }
  }

  @action
  addTarget() {
    const { label, key } = this.selectedTarget;
    const value = this.selectedTargetValue;
    this.targets = addToArray(this.targets, { label, value, key });
    // recalculate value for appropriate model property
    this.updateModelForKey(key);
    this.selectedTargetValue = null;
    this.resetTargetState();
  }
  @action
  removeTarget(target) {
    this.targets = removeFromArray(this.targets, target);
    // recalculate value for appropriate model property
    this.updateModelForKey(target.key);
  }
  @action
  cancel() {
    // revert model changes
    this.args.onClose();
  }
}
