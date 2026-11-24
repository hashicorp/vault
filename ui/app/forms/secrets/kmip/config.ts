/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import OpenApiForm from 'vault/forms/open-api';

import type { KmipConfigureRequest } from '@hashicorp/vault-client-typescript';
import type Form from 'vault/forms/form';
import type FormField from 'vault/utils/forms/field';

export default class KmipConfigForm extends OpenApiForm<KmipConfigureRequest> {
  constructor(...args: ConstructorParameters<typeof Form>) {
    super('KmipConfigureRequest', ...args);

    const orderedKeys = [
      'listen_addrs',
      'default_tls_client_key_bits',
      'default_tls_client_key_type',
      'default_tls_client_ttl',
      'server_hostnames',
      'server_ips',
      'tls_ca_key_bits',
      'tls_ca_key_type',
      'tls_min_version',
    ];

    this.formFields = orderedKeys.map((key) => this.formFields.find((f) => f.name === key) as FormField);
  }
}
