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
import SshSignForm from 'vault/forms/ssh/sign';

import type ApiService from 'vault/services/api';
import type ControlGroupService from 'vault/vault/services/control-group';

interface Args {
  roleName: string;
  backendPath: string;
}

export default class SshSignKey extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly controlGroup: ControlGroupService;

  @tracked signForm = new SshSignForm({ cert_type: 'user' });
  @tracked signedKeyData: Record<string, unknown> | null = null;
  @tracked errorMessage: string | null = null;
  @tracked modelValidations: Record<string, unknown> | null = null;
  @tracked invalidFormAlert: string | null = null;

  get signDisplayRows() {
    const data = this.signedKeyData;
    if (!data) return [];
    return [
      { label: 'Signed key', value: data['signed_key'] },
      { label: 'Lease ID', value: data['lease_id'] },
      { label: 'Renewable', value: data['renewable'] },
      { label: 'Lease duration', value: data['lease_duration'] },
      { label: 'Serial number', value: data['serial_number'] },
    ].filter((f) => f.value != null && f.value !== '');
  }

  get breadcrumbs() {
    const { backendPath, roleName } = this.args;
    return [
      { label: backendPath, route: 'vault.cluster.secrets.backend', model: backendPath },
      { label: 'Sign', route: 'vault.cluster.secrets.backend', model: backendPath },
      { label: roleName, route: 'vault.cluster.secrets.backend.show', model: roleName },
      { label: 'Sign SSH Key' },
    ];
  }

  sign = task(
    waitFor(async (evt: Event) => {
      evt.preventDefault();
      this.errorMessage = null;

      const { isValid, state, invalidFormMessage, data } = this.signForm.toJSON();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = isValid ? null : invalidFormMessage;
      if (!isValid) return;
      try {
        const result = await this.api.secrets.sshSignCertificate(
          this.args.roleName,
          this.args.backendPath,
          data
        );
        this.signedKeyData = { ...result, ...(result.data as Record<string, unknown>) };
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
    this.signedKeyData = null;
    this.signForm = new SshSignForm({ cert_type: 'user' });
    this.errorMessage = null;
    this.modelValidations = null;
    this.invalidFormAlert = null;
  }
}
