/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type SecretMountPath from 'vault/services/secret-mount-path';
import { KmipGenerateClientCertificateRequest } from '@hashicorp/vault-client-typescript';
import { HTMLElementEvent } from 'vault/forms';

type Credentials = {
  certificate: string;
  serial_number: string;
  ca_chain: string[];
  private_key: string;
};

interface Args {
  scopeName: string;
  roleName: string;
}

export default class KmipPageCredentialsGenerateComponent extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly secretMountPath: SecretMountPath;

  @tracked declare credentials: Credentials;
  @tracked format = 'pem';
  @tracked invalidFormAlert: string | null = null;
  @tracked errorMessage: string | null = null;

  generate = task(
    waitFor(async (event: HTMLElementEvent<HTMLFormElement>) => {
      event.preventDefault();
      const { scopeName, roleName } = this.args;
      const { currentPath } = this.secretMountPath;

      try {
        const payload = { format: this.format } as KmipGenerateClientCertificateRequest;
        const { data } = await this.api.secrets.kmipGenerateClientCertificate(
          roleName,
          scopeName,
          currentPath,
          payload
        );
        this.credentials = data as Credentials;
        this.flashMessages.success(`Successfully generated credentials from role ${roleName}.`);
      } catch (error) {
        const { message } = await this.api.parseError(error);
        this.errorMessage = message;
        this.invalidFormAlert = 'There was an error submitting this form.';
      }
    })
  );
}
