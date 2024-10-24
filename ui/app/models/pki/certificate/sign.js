/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import PkiCertificateBaseModel from './base';

const generateFromRole = [
  {
    default: ['csr', 'commonName', 'customTtl', 'format', 'removeRootsFromChain'],
  },
  {
    'Subject Alternative Name (SAN) Options': [
      'excludeCnFromSans',
      'altNames',
      'ipSans',
      'uriSans',
      'otherSans',
    ],
  },
];
@withFormFields(null, generateFromRole)
export default class PkiCertificateSignModel extends PkiCertificateBaseModel {
  @attr('string') role; // role name to create certificate against for request URL

  @attr('string', {
    label: 'CSR',
    editType: 'textarea',
  })
  csr;

  @attr('boolean', {
    subText: 'When checked, the CA chain will not include self-signed CA certificates.',
  })
  removeRootsFromChain;
}
