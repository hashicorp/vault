/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import OpenApiForm from 'vault/forms/open-api';
import FormField from 'vault/utils/forms/field';

import type {
  PkiGenerateRootRequest,
  PkiIssuersGenerateRootRequest,
  PkiGenerateIntermediateRequest,
  PkiIssuersGenerateIntermediateRequest,
} from '@hashicorp/vault-client-typescript';
import type { Validations } from 'vault/app-types';

type PkiGenerateRequest =
  | PkiGenerateRootRequest
  | PkiIssuersGenerateRootRequest
  | PkiGenerateIntermediateRequest
  | PkiIssuersGenerateIntermediateRequest;
type PkiConfigGenerateFormData = PkiGenerateRequest & {
  type?: string;
  customTtl?: string;
};

export default class PkiConfigGenerateForm extends OpenApiForm<PkiConfigGenerateFormData> {
  constructor(...args: ConstructorParameters<typeof OpenApiForm>) {
    super(...args);
    // type and customTtl are UI only fields used to determine which fields to show and which validations to apply
    // add them manually to the formFields since they are not included in the helpResponse
    this.formFields.push(
      new FormField('type', 'string', {
        possibleValues: ['exported', 'internal', 'existing', 'kms'],
        noDefault: true,
      }),
      new FormField('customTtl', 'string', {
        label: 'Not valid after',
        subText:
          'The time after which this certificate will no longer be valid. This can be a TTL (a range of time from now) or a specific date.',
        editType: 'yield',
      })
    );
  }

  validations: Validations = {
    type: [{ type: 'presence', message: 'Type is required.' }],
    common_name: [{ type: 'presence', message: 'Common name is required.' }],
    issuer_name: [
      {
        type: 'isNot',
        options: { value: 'default' },
        message: `Issuer name must be unique across all issuers and not be the reserved value 'default'.`,
      },
    ],
    key_name: [
      {
        type: 'isNot',
        options: { value: 'default' },
        message: `Key name cannot be the reserved value 'default'`,
      },
    ],
  };

  toJSON() {
    const { data, ...rest } = super.toJSON(this.data);
    // remove type and customTtl from payload since they are UI only props
    const { type, customTtl, ...payload } = data;
    return { data: payload, ...rest };
  }
}
