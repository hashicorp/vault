/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { addManyToArray, addToArray } from 'vault/helpers/add-to-array';
import { removeFromArray } from 'vault/helpers/remove-from-array';

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
 * @param {Object} model - login enforcement model
 * @param {Object} [isInline] - toggles inline display of form -- method selector and actions are hidden and should be handled externally
 * @param {Object} [modelErrors] - model validations state object if handling actions externally when displaying inline
 * @param {onSave} [onSave] - triggered on save success
 * @param {onClose} [onClose] - triggered on cancel
 */

export default class MfaLoginEnforcementForm extends Component {
  @service store;
  @service flashMessages;

  targetTypes = [
    { label: 'Authentication mount', type: 'accessor', key: 'auth_method_accessors' },
    { label: 'Authentication method', type: 'method', key: 'auth_method_types' },
    { label: 'Group', type: 'identity/group', key: 'identity_groups' },
    { label: 'Entity', type: 'identity/entity', key: 'identity_entities' },
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

  constructor() {
    super(...arguments);
    // aggregate different target array properties on model into flat list
    this.flattenTargets();
    // eagerly fetch identity groups and entities for use as search select options
    this.resetTargetState();
    // only auth method types that have mounts can be selected as targets -- fetch from sys/auth and map by type
    this.fetchAuthMethods();
  }

  async flattenTargets() {
    for (const { label, key } of this.targetTypes) {
      const targetArray = await this.args.model[key];
      const targets = targetArray.map((value) => ({ label, key, value }));
      this.targets = addManyToArray(this.targets, targets);
    }
  }
  async resetTargetState() {
    this.selectedTargetValue = null;
    const options = this.searchSelectOptions || {};
    if (!this.searchSelectOptions) {
      const types = ['identity/group', 'identity/entity'];
      for (const type of types) {
        try {
          options[type] = await this.store.query(type, {});
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
    const mounts = await this.store.findAll('auth-method');
    this.authMethods = mounts.map((auth) => auth.type);
  }

  get selectedTarget() {
    return this.targetTypes.find((tt) => tt.type === this.selectedTargetType);
  }
  get errors() {
    return this.args.modelErrors || this.modelErrors;
  }

  updateModelForKey(key) {
    const newValue = this.targets.filter((t) => t.key === key).map((t) => t.value);
    this.args.model[key] = newValue;
  }

  @task
  *save() {
    this.modelErrors = {};
    // check validity state first and abort if invalid
    const { isValid, state } = this.args.model.validate();
    if (!isValid) {
      this.modelErrors = state;
    } else {
      try {
        yield this.args.model.save();
        this.args.onSave();
      } catch (error) {
        const message = error.errors ? error.errors.join('. ') : error.message;
        this.flashMessages.danger(message);
      }
    }
  }

  @action
  async onMethodChange(selectedIds) {
    // first make sure the async relationship is loaded
    const methods = await this.args.model.mfa_methods;
    // then remove items that are no longer selected
    const updatedList = methods.filter((model) => {
      return selectedIds.includes(model.id);
    });
    // then add selected items that don't exist in the list already
    const modelIds = updatedList.map((model) => model.id);
    const toAdd = selectedIds
      .filter((id) => !modelIds.includes(id))
      .map((id) => this.store.peekRecord('mfa-method', id));
    this.args.model.mfa_methods = addManyToArray(updatedList, toAdd);
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
      // for identity groups and entities grab model from store as value
      this.selectedTargetValue = this.store.peekRecord(type, selected[0]);
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
    this.args.model.rollbackAttributes();
    this.args.onClose();
  }
}
