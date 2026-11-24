/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import SshOtpCredentialForm from 'vault/forms/ssh/otp-credential';

import type ApiService from 'vault/services/api';
import type ControlGroupService from 'vault/vault/services/control-group';

interface Args {
  backendPath: string;
  roleName: string;
}

export default class GenerateCredentialsSsh extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly controlGroup: ControlGroupService;

  @tracked credentialForm = new SshOtpCredentialForm();
  @tracked otpData: Record<string, unknown> | null = null;
  @tracked errorMessage: string | null = null;
  @tracked modelValidations: Record<string, unknown> | null = null;
  @tracked invalidFormAlert: string | null = null;

  get otpDisplayRows() {
    const data = this.otpData;
    if (!data) return [];
    return [
      { label: 'Username', value: data['username'] },
      { label: 'IP Address', value: data['ip'] },
      { label: 'Key', value: data['key'], masked: true },
      { label: 'Key type', value: data['key_type'] },
      { label: 'Port', value: data['port'] },
    ].filter((f) => f.value != null && f.value !== '');
  }

  get breadcrumbs() {
    const { backendPath, roleName } = this.args;
    return [
      { label: backendPath, route: 'vault.cluster.secrets.backend', model: backendPath },
      { label: 'Credentials', route: 'vault.cluster.secrets.backend', model: backendPath },
      { label: roleName, route: 'vault.cluster.secrets.backend.show', model: roleName },
      { label: 'Generate SSH credentials' },
    ];
  }

  generate = task(
    waitFor(async (evt: Event) => {
      evt.preventDefault();
      this.errorMessage = null;

      const { isValid, state, invalidFormMessage, data } = this.credentialForm.toJSON();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = isValid ? null : invalidFormMessage;
      if (!isValid) return;
      try {
        const result = await this.api.secrets.sshGenerateCredentials(
          this.args.roleName,
          this.args.backendPath,
          data
        );
        this.otpData = (result.data as Record<string, unknown>) ?? {};
      } catch (error) {
        const { message, response } = await this.api.parseError(error);
        if (response?.isControlGroupError) {
          this.controlGroup.saveTokenFromError(response);
          this.errorMessage = this.controlGroup.logFromError(response).content;
        } else {
          this.errorMessage = message;
        }
      }
    })
  );

  @action
  reset() {
    this.otpData = null;
    this.credentialForm = new SshOtpCredentialForm();
    this.errorMessage = null;
    this.modelValidations = null;
    this.invalidFormAlert = null;
  }
}
