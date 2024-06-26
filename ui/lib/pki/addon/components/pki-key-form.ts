/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';
import type FlashMessageService from 'vault/services/flash-messages';
import type PkiKeyModel from 'vault/models/pki/key';
import type { ValidationMap } from 'vault/app-types';

/**
 * @module PkiKeyForm
 * PkiKeyForm components are used to create and update PKI keys.
 *
 * @example
 * ```js
 * <PkiKeyForm @model={{this.model}} @onCancel={{transition-to "vault.cluster"}} @onSave={{transition-to "vault.cluster"}} />
 * ```
 *
 * @param {Object} model - pki/key model.
 * @callback onCancel - Callback triggered when cancel button is clicked.
 * @callback onSave - Callback triggered on save success.
 */

interface Args {
  model: PkiKeyModel;
  onSave: CallableFunction;
}

export default class PkiKeyForm extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;

  @tracked errorBanner = '';
  @tracked invalidFormAlert = '';
  @tracked modelValidations: ValidationMap | null = null;

  @task
  @waitFor
  *save(event: Event) {
    event.preventDefault();
    try {
      const { isNew, keyName } = this.args.model;
      const { isValid, state, invalidFormMessage } = this.args.model.validate();
      if (isNew) {
        this.modelValidations = isValid ? null : state;
        this.invalidFormAlert = invalidFormMessage;
      }
      if (!isValid && isNew) return;
      yield this.args.model.save({ adapterOptions: { import: false } });
      this.flashMessages.success(
        `Successfully ${isNew ? 'generated' : 'updated'} key${keyName ? ` ${keyName}.` : '.'}`
      );
      this.args.onSave();
    } catch (error) {
      this.errorBanner = errorMessage(error);
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }
}
