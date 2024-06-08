/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint ember/no-computed-properties-in-native-classes: 'warn' */
/**
 * @module TotpEdit
 * TotpEdit component manages the secret and model data, and displays either the create, empty state or show view of a TOTP account.
 *
 * @example
 * ```js
 * <TotpEdit @model={{model}} @mode="create" @baseKey={{this.baseKey}} @key={{this.model}} @initialKey={{this.initialKey}} @onRefresh={{action "refresh"}} @onToggleAdvancedEdit={{action "toggleAdvancedEdit"}} @preferAdvancedEdit={{this.preferAdvancedEdit}}/>
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
import { later, cancel } from '@ember/runloop';
import timestamp from 'core/utils/timestamp';
import { getUnixTime } from 'date-fns';

export default class TotpEdit extends Component {
  @service store;

  @tracked elapsedTime = 0;
  nextTick = null;

  get remainingTime() {
    const { model } = this.args;
    return model.period - this.elapsedTime;
  }

  @action
  cancelTimer() {
    cancel(this.nextTick);
  }

  @action
  startTimer() {
    this.nextTick = later(
      this,
      function () {
        const { model } = this.args;
        this.elapsedTime = getUnixTime(timestamp.now()) % model.period;
        this.startTimer();
      },
      1000
    );
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
  @alias('checkSecretCapabilities.canRead') canReadSecret;

  @or('model.isLoading', 'model.isReloading', 'model.isSaving') requestInFlight;
  @or('requestInFlight', 'model.isFolder', 'model.flagsIsInvalid') buttonDisabled;

  get modelForData() {
    const { model } = this.args;
    if (!model) return null;
    return model;
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
}
