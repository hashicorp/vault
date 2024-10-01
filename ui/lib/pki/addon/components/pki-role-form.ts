/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import errorMessage from 'vault/utils/error-message';
import type Store from '@ember-data/store';
import type FlashMessageService from 'vault/services/flash-messages';
import type SecretMountPathService from 'vault/services/secret-mount-path';
import type PkiRoleModel from 'vault/models/pki/role';
import type PkiIssuerModel from 'vault/models/pki/issuer';
import type { ValidationMap } from 'vault/app-types';

/**
 * @module PkiRoleForm
 * PkiRoleForm components are used to create and update PKI roles.
 *
 * @example
 * ```js
 * <PkiRoleForm @model={{this.model}}/>
 * ```
 * @callback onCancel
 * @callback onSave
 * @param {Object} role - pki/role model.
 * @param {Array} issuers - pki/issuer model.
 * @param {onCancel} onCancel - Callback triggered when cancel button is clicked.
 * @param {onSave} onSave - Callback triggered on save success.
 */

interface Args {
  role: PkiRoleModel;
  issuers: PkiIssuerModel[];
  onSave: CallableFunction;
}

export default class PkiRoleForm extends Component<Args> {
  @service declare readonly store: Store;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly secretMountPath: SecretMountPathService;

  @tracked errorBanner = '';
  @tracked invalidFormAlert = '';
  @tracked modelValidations: ValidationMap | null = null;
  @tracked showDefaultIssuer = true;

  constructor(owner: unknown, args: Args) {
    super(owner, args);

    this.showDefaultIssuer = this.args.role.issuerRef === 'default';
  }

  get issuers() {
    return this.args.issuers?.map((issuer) => {
      return { issuerDisplayName: issuer.issuerName || issuer.issuerId };
    });
  }

  @task
  *save(event: Event) {
    event.preventDefault();
    try {
      const { isValid, state, invalidFormMessage } = this.args.role.validate();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = invalidFormMessage;
      if (isValid) {
        const { isNew, name } = this.args.role;
        yield this.args.role.save();
        this.flashMessages.success(`Successfully ${isNew ? 'created' : 'updated'} the role ${name}.`);
        this.args.onSave();
      }
    } catch (error) {
      this.errorBanner = errorMessage(error);
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }

  @action
  toggleShowDefaultIssuer() {
    this.showDefaultIssuer = !this.showDefaultIssuer;

    if (this.showDefaultIssuer) {
      this.args.role.issuerRef = 'default';
    }
  }
}
