/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import OpenApiForm from 'vault/forms/open-api';
import FormField from 'vault/utils/forms/field';

import type { PkiIssuerSignIntermediateRequest } from '@hashicorp/vault-client-typescript';
import type Form from 'vault/forms/form';
import type { Validations } from 'vault/vault/app-types';

export default class PkiIssuersSignIntermediateForm extends OpenApiForm<PkiIssuerSignIntermediateRequest> {
  constructor(...args: ConstructorParameters<typeof Form>) {
    super('PkiIssuerSignIntermediateRequest', ...args);
    // customTtl is a convenience field that sets ttl and notAfter via one input <PkiNotValidAfterForm>
    this.formFields.push(
      new FormField('customTtl', undefined, {
        label: 'Not valid after',
        subText:
          'The time after which this certificate will no longer be valid. This can be a TTL (a range of time from now) or a specific date.',
        editType: 'yield',
      })
    );
    // better UX for csr field to be textarea
    const csrField = this.formFields.find((f) => f.name === 'csr');
    if (csrField) {
      csrField.options.editType = 'textarea';
    }
  }

  validations: Validations = {
    csr: [{ type: 'presence', message: 'CSR is required.' }],
  };
}
