/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import PkiCertificateBaseModel from './base';

const generateFromRole = [
  {
    default: ['commonName'],
  },
  {
    'Subject Alternative Name (SAN) Options': [
      'altNames',
      'ipSans',
      'uriSans',
      'otherSans',
      'ttl',
      'format',
      'privateKeyFormat',
      'excludeCnFromSans',
      'notAfter',
    ],
  },
];
@withFormFields(null, generateFromRole)
export default class PkiCertificateGenerateModel extends PkiCertificateBaseModel {
  getHelpUrl(backend) {
    return `/v1/${backend}/issue/example?help=1`;
  }
  @attr('string') role;
}
