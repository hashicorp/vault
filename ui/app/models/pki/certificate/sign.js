/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
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
  getHelpUrl(backend) {
    return `/v1/${backend}/sign/example?help=1`;
  }
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
