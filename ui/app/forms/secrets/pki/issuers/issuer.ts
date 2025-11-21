/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import OpenApiForm from 'vault/forms/open-api';

import type { PkiWriteIssuerRequest } from '@hashicorp/vault-client-typescript';
import type Form from 'vault/forms/form';

export default class PkiIssuerForm extends OpenApiForm<PkiWriteIssuerRequest> {
  constructor(...args: ConstructorParameters<typeof Form>) {
    super('PkiWriteIssuerRequest', ...args);

    this.formFields.forEach((field) => {
      // usage and leaf_not_after_behavior require special handling in the UI - update editType to yield
      if (['leaf_not_after_behavior', 'usage'].includes(field.name)) {
        field.options.editType = 'yield';
        // the form field component conditionals should be reworked to yield no matter the type
        // for now set type to undefined to ensure block is yielded out
        field.type = undefined;
      }
      // add options for revocation_signature_algorithm
      if (field.name === 'revocation_signature_algorithm') {
        field.options.noDefault = true;
        field.options.possibleValues = [
          'SHA256WithRSA',
          'ECDSAWithSHA384',
          'SHA256WithRSAPSS',
          'ED25519',
          'SHA384WithRSAPSS',
          'SHA512WithRSAPSS',
          'PureEd25519',
          'SHA384WithRSA',
          'SHA512WithRSA',
          'ECDSAWithSHA256',
          'ECDSAWithSHA512',
        ];
      }
    });
  }
}
