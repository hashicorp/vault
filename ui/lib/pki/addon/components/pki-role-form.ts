/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type SecretMountPathService from 'vault/services/secret-mount-path';
import type VersionService from 'vault/services/version';
import type PkiRoleForm from 'vault/forms/secrets/pki/role';
import type { ValidationMap } from 'vault/app-types';

/**
 * @module PkiRoleFormComponent
 * PkiRoleFormComponent is used to create and update PKI roles.
 *
 * @callback onCancel
 * @callback onSave
 * @param {Object} role - pki role form class.
 * @param {Array} issuers - pki issuers list key info.
 * @param {onCancel} onCancel - Callback triggered when cancel button is clicked.
 * @param {onSave} onSave - Callback triggered on save success.
 */

interface Args {
  form: PkiRoleForm;
  issuers: { issuer_name?: string; issuer_id: string }[];
  onSave: CallableFunction;
  onCancel: CallableFunction;
}

export default class PkiRoleFormComponent extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly secretMountPath: SecretMountPathService;
  @service declare readonly version: VersionService;

  @tracked errorBanner = '';
  @tracked invalidFormAlert = '';
  @tracked modelValidations: ValidationMap | null = null;
  @tracked showDefaultIssuer = true;
  @tracked openGroups: string[] = [];

  constructor(owner: unknown, args: Args) {
    super(owner, args);

    this.showDefaultIssuer = this.args.form.data.issuer_ref === 'default';
  }

  // hide no_store_metadata field for community edition
  showField = (fieldName: string) => (fieldName === 'no_store_metadata' ? this.version.isEnterprise : true);

  get issuers() {
    return this.args.issuers?.map(({ issuer_name, issuer_id }) => {
      return { issuerDisplayName: issuer_name || issuer_id };
    });
  }

  save = task(
    waitFor(async (event: Event) => {
      event.preventDefault();
      try {
        const { isValid, state, invalidFormMessage, data } = this.args.form.toJSON();
        this.modelValidations = isValid ? null : state;
        this.invalidFormAlert = invalidFormMessage;
        if (isValid) {
          const { name, ...payload } = data;
          await this.api.secrets.pkiWriteRole(name, this.secretMountPath.currentPath, payload);
          this.flashMessages.success(`Successfully saved the role ${name}.`);
          this.args.onSave(name);
        }
      } catch (error) {
        const { message } = await this.api.parseError(error);
        this.errorBanner = message;
        this.invalidFormAlert = 'There was an error submitting this form.';
      }
    })
  );

  @action
  toggleShowDefaultIssuer() {
    this.showDefaultIssuer = !this.showDefaultIssuer;

    if (this.showDefaultIssuer) {
      this.args.form.data.issuer_ref = 'default';
    }
  }

  @action
  toggleGroup(group: string) {
    if (this.openGroups.includes(group)) {
      this.openGroups = this.openGroups.filter((g) => g !== group);
    } else {
      this.openGroups = [...this.openGroups, group];
    }
  }
}
