/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import PkiCertificateBaseModel from './base';

const generateFromRole = [
  {
    default: ['commonName', 'userIds', 'customTtl', 'format', 'privateKeyFormat'],
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
// Extra fields returned on the /issue endpoint
const certDisplayFields = [
  'certificate',
  'commonName',
  'revocationTime',
  'serialNumber',
  'caChain',
  'issuingCa',
  'privateKey',
  'privateKeyType',
];
@withFormFields(certDisplayFields, generateFromRole)
export default class PkiCertificateGenerateModel extends PkiCertificateBaseModel {
  getHelpUrl(backend) {
    return `/v1/${backend}/issue/example?help=1`;
  }
  @attr('string') role; // role name to issue certificate against for request URL
}
