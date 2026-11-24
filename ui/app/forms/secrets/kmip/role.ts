/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import OpenApiForm from 'vault/forms/open-api';

import type { KmipWriteRoleRequest } from '@hashicorp/vault-client-typescript';
import type Form from 'vault/forms/form';
import type FormField from 'vault/utils/forms/field';

export default class KmipRoleForm extends OpenApiForm<KmipWriteRoleRequest> {
  constructor(...args: ConstructorParameters<typeof Form>) {
    super('KmipWriteRoleRequest', ...args);
  }

  get tlsFields() {
    return this.formFields.filter((field) => field.name.startsWith('tls_'));
  }
  // there are currently no other fields but adding this for future-proofing
  get otherFields() {
    return this.formFields.filter(
      (field) => !field.name.startsWith('tls_') && !field.name.startsWith('operation_')
    );
  }
  // helper used in form template to look up field by name
  fieldFor = (key: keyof KmipWriteRoleRequest) => {
    return this.formFields.find((f) => f.name === key) as FormField;
  };

  toJSON() {
    let data = this.data;
    const { tls_client_key_bits, tls_client_key_type, tls_client_ttl } = data;
    const tls = { tls_client_key_bits, tls_client_key_type, tls_client_ttl };
    if (data.operation_all) {
      data = { ...tls, operation_all: true };
    } else if (data.operation_none) {
      data = { ...tls, operation_none: true };
    } else {
      // ensure operation_all and operation_none are not present in payload
      const { operation_all, operation_none, ...rest } = data;
      data = rest;
    }
    return super.toJSON(data);
  }
}
