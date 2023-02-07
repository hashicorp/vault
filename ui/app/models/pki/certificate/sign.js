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
  @attr('string') role;

  @attr('string', {
    label: 'CSR',
    editType: 'textarea',
  })
  csr;

  @attr({
    label: 'Not valid after',
    detailsLabel: 'Issued certificates expire after',
    subText:
      'The time after which this certificate will no longer be valid. This can be a TTL (a range of time from now) or a specific date.',
    editType: 'yield',
  })
  customTtl;

  @attr('boolean', {
    subText: 'When checked, the CA chain will not include self-signed CA certificates',
  })
  removeRootsFromChain;
}
