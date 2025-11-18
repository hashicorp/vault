/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { toLabel } from 'core/helpers/to-label';

import type { PkiReadIssuerResponse } from '@hashicorp/vault-client-typescript';
import type { ParsedCertificateData } from 'vault/utils/parse-pki-cert';

interface Args {
  issuer: PkiReadIssuerResponse & { parsedCertificate: ParsedCertificateData; isRoot: boolean };
  pem: string;
  der: string;
  isRotatable: boolean;
  canRotate: boolean;
  canCrossSign: boolean;
  canSignIntermediate: boolean;
  canConfigure: boolean;
  backend: string;
}

export default class PkiIssuerDetailsComponent extends Component<Args> {
  @tracked showRotationModal = false;

  defaultFields = [
    'certificate',
    'ca_chain',
    'parsedCertificate.common_name',
    'issuer_name',
    'issuer_id',
    'key_id',
  ];
  urlFields = ['issuing_certificates_urls', 'crl_distribution_points', 'ocsp_servers'];

  label = (field: string) => {
    const label = toLabel([field]);
    return (
      {
        ca_chain: 'CA Chain',
        'parsedCertificate.common_name': 'Common name',
        issuer_id: 'Issuer ID',
        key_id: 'Default key ID',
        crl_distribution_points: 'CRL distribution points',
        ocsp_servers: 'OCSP servers',
      }[field] || label
    );
  };

  get parsingErrors() {
    const { parsedCertificate } = this.args.issuer;
    if (parsedCertificate?.parsing_errors?.length) {
      return parsedCertificate.parsing_errors.map((e: Error) => e.message).join(', ');
    }
    return '';
  }
}
