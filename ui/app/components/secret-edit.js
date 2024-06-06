/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint ember/no-computed-properties-in-native-classes: 'warn' */
/**
 * @module SecretEdit
 * SecretEdit component manages the secret and model data, and displays either the create, update, empty state or show view of a KV secret.
 *
 * @example
 * ```js
 * <SecretEdit @model={{model}} @mode="create" @baseKey={{this.baseKey}} @key={{this.model}} @initialKey={{this.initialKey}} @onRefresh={{action "refresh"}} @onToggleAdvancedEdit={{action "toggleAdvancedEdit"}} @preferAdvancedEdit={{this.preferAdvancedEdit}}/>
 * ```
/This component is initialized from the secret-edit-layout.hbs file
 * @param {object} model - Secret model which is generated in the secret-edit route
 * @param {string} mode - Edit, create, etc.
 * @param {string} baseKey - Provided for navigation.
 * @param {object} key - Passed through, copy of the model.
 * @param {string} initialKey - model's name.
 * @param {function} onRefresh - action that refreshes the model
 * @param {function} onToggleAdvancedEdit - changes the preferAdvancedEdit to true or false
 * @param {boolean} preferAdvancedEdit - property set from the controller of show/edit/create route passed in through secret-edit-layout
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import KVObject from 'vault/lib/kv-object';
import { maybeQueryRecord } from 'vault/macros/maybe-query-record';
import { alias, or } from '@ember/object/computed';

export default class SecretEdit extends Component {
  @service store;

  @tracked secretData = null;
  @tracked codemirrorString = null;

  // fired on did-insert from render modifier
  @action
  createKvData(elem, [model]) {
    this.secretData = KVObject.create({ content: [] }).fromJSON(model.secretData);
    this.codemirrorString = this.secretData.toJSONString();
  }
  // TODO move this to the secret model
  @maybeQueryRecord(
    'capabilities',
    (context) => {
      if (!context.args.model || context.args.mode === 'create') {
        return;
      }
      const backend = context.args.model.backend;
      const id = context.args.model.id;
      const path = `${backend}/${id}`;
      return {
        id: path,
      };
    },
    'model',
    'model.id',
    'mode'
  )
  checkSecretCapabilities;
  @alias('checkSecretCapabilities.canUpdate') canUpdateSecret;
  @alias('checkSecretCapabilities.canRead') canReadSecret;

  @or('model.isLoading', 'model.isReloading', 'model.isSaving') requestInFlight;
  @or('requestInFlight', 'model.isFolder', 'model.flagsIsInvalid') buttonDisabled;

  get modelForData() {
    const { model } = this.args;
    if (!model) return null;
    return model;
  }

  get basicModeDisabled() {
    return this.secretDataIsAdvanced || this.showAdvancedMode === false;
  }

  get secretDataAsJSON() {
    return this.secretData.toJSON();
  }

  get secretDataIsAdvanced() {
    return this.secretData.isAdvanced();
  }

  get showAdvancedMode() {
    return this.secretDataIsAdvanced || this.args.preferAdvancedEdit;
  }

  get isWriteWithoutRead() {
    if (!this.args.model) {
      return false;
    }
    // if the model couldn't be read from the server
    if (this.args.model.failedServerRead) {
      return true;
    }
    return false;
  }

  @action
  refresh() {
    this.args.onRefresh();
  }

  @action
  toggleAdvanced(bool) {
    this.args.onToggleAdvancedEdit(bool);
  }
}
