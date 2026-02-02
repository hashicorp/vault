/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import OpenApiForm from 'vault/forms/open-api';
import FormFieldGroup from 'vault/utils/forms/field-group';
import FormField from 'vault/utils/forms/field';

import type { PkiIssueWithRoleRequest, PkiSignWithRoleRequest } from '@hashicorp/vault-client-typescript';

type PkiCertificateFormData = PkiIssueWithRoleRequest | PkiSignWithRoleRequest;

export default class PkiCertificateForm extends OpenApiForm<PkiCertificateFormData> {
  constructor(...args: ConstructorParameters<typeof OpenApiForm>) {
    super(...args);

    const sansKeys = ['exclude_cn_from_sans', 'alt_names', 'ip_sans', 'uri_sans', 'other_sans'];
    const excludeKeys = ['ttl', 'issuer_ref']; // issuer_ref is not editable, ttl is set via customTtl
    const primaryKeys = ['common_name', 'csr'];
    const defaultGroup = this.formFieldGroups[0]?.['default'] as FormField[];

    const fields = defaultGroup.reduce(
      (fields: { default: FormField[]; sans: FormField[] }, field: FormField) => {
        // better UX if csr is textarea
        if (field.name === 'csr') {
          field.options.editType = 'textarea';
        }
        // move sans related fields to their own group
        if (sansKeys.includes(field.name)) {
          fields.sans.push(field);
        } else if (field.name === 'not_after') {
          // customTtl is a convenience field that sets ttl and notAfter via one input <PkiNotValidAfterForm>
          // remove not_after and ttl fields and replace with customTtl
          const customTtlField = new FormField('customTtl', undefined, {
            label: 'Not valid after',
            subText:
              'The time after which this certificate will no longer be valid. This can be a TTL (a range of time from now) or a specific date.',
            editType: 'yield',
          });
          fields.default.push(customTtlField);
        } else if (primaryKeys.includes(field.name)) {
          // move common_name and csr (sign only) to the top of the default group
          fields.default.unshift(field);
        } else if (!excludeKeys.includes(field.name)) {
          fields.default.push(field);
        }
        return fields;
      },
      { default: [], sans: [] }
    );

    this.formFieldGroups = [
      new FormFieldGroup('default', fields.default),
      new FormFieldGroup('Subject Alternative Name (SAN) Options', fields.sans),
    ];
  }
}
