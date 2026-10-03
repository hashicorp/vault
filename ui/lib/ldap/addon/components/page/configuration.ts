/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { toLabel } from 'core/helpers/to-label';

import type FlashMessageService from 'vault/services/flash-messages';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type { Breadcrumb } from 'vault/vault/app-types';
import type { LdapApplicationModel } from 'ldap/routes/application';

interface Args {
  model: LdapApplicationModel;
  breadcrumbs: Array<Breadcrumb>;
}

export default class LdapConfigurationPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  defaultFields = [
    'binddn',
    'url',
    'schema',
    'password_policy',
    'userdn',
    'userattr',
    'connection_timeout',
    'request_timeout',
  ];

  connectionFields = ['certificate', 'starttls', 'insecure_tls', 'client_tls_cert', 'client_tls_key'];

  label = (field: string) => {
    return (
      {
        binddn: 'Administrator distinguished name',
        url: 'URL',
        certificate: 'CA certificate',
        starttls: 'Start TLS',
        insecure_tls: 'Insecure TLS',
        client_tls_cert: 'Client TLS certificate',
        client_tls_key: 'Client TLS key',
      }[field] || toLabel([field])
    );
  };

  rotateRoot = task(
    waitFor(async () => {
      try {
        await this.api.secrets.ldapRotateRootCredentials(this.secretMountPath.currentPath);
        this.flashMessages.success('Root password successfully rotated.');
      } catch (error) {
        const { message } = await this.api.parseError(error);
        this.flashMessages.danger(`Error rotating root password \n ${message}`);
      }
    })
  );
}
