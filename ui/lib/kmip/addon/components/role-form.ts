/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

import type KmipRoleForm from 'vault/forms/secrets/kmip/role';
import type FlashMessageService from 'vault/services/flash-messages';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type { HTMLElementEvent } from 'vault/forms';

interface Args {
  roleName: string;
  scopeName: string;
  form: KmipRoleForm;
  onSave: CallableFunction;
  onCancel: CallableFunction;
}

export default class KmipRoleFormComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  @tracked declare name: string;
  @tracked validationError = false;
  @tracked invalidFormAlert: string | null = null;
  @tracked errorMessage: string | null = null;

  @action
  toggleOperationNone(event: HTMLElementEvent<HTMLInputElement>) {
    const { checked } = event.target;
    const { data } = this.args.form;
    data.operation_none = !checked;
    data.operation_all = checked;
  }

  save = task(
    waitFor(async (event: HTMLElementEvent<HTMLFormElement>) => {
      event.preventDefault();
      const { form, roleName, scopeName } = this.args;

      if (!form.isNew || this.name) {
        this.validationError = false;
        try {
          const { data } = form.toJSON();
          const name = form.isNew ? this.name : roleName;

          await this.api.secrets.kmipWriteRole(name, scopeName, this.secretMountPath.currentPath, data);

          this.flashMessages.success(`Successfully saved role ${name}`);
          this.args.onSave();
        } catch (error) {
          const { message } = await this.api.parseError(error);
          this.errorMessage = message;
          this.invalidFormAlert = 'There was an error submitting this form.';
        }
      } else {
        this.validationError = true;
        this.invalidFormAlert = 'There is an error with this form.';
      }
    })
  );
}
