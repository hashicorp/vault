/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { toLabel } from 'core/helpers/to-label';
import errorMessage from 'vault/utils/error-message';
import timestamp from 'core/utils/timestamp';

import type DownloadService from 'vault/services/download';
import type { Extensions } from 'vault/services/download';
import type FlashMessageService from 'vault/services/flash-messages';
import type { PkiReadIssuerResponse } from '@hashicorp/vault-client-typescript';
import type { ParsedCertificateData } from 'vault/utils/parse-pki-cert';

interface Args {
  issuer: PkiReadIssuerResponse & { parsedCertificate: ParsedCertificateData; isRoot: boolean };
  pem: string;
  der: Blob;
  isRotatable: boolean;
  canRotate: boolean;
  canCrossSign: boolean;
  canSignIntermediate: boolean;
  canConfigure: boolean;
  backend: string;
}

export default class PkiIssuerDetailsComponent extends Component<Args> {
  @service declare readonly download: DownloadService;
  @service declare readonly flashMessages: FlashMessageService;

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

  // Replaces <DownloadButton> so the menu rows can be <Hds::Dropdown> list items.
  @action
  downloadIssuer(format: 'der' | 'pem', close: CallableFunction) {
    const { issuer_id: issuerId } = this.args.issuer;
    const content = format === 'der' ? this.args.der : this.args.pem;
    // Matches the filename <DownloadButton> produces.
    const ts = timestamp.now().toISOString();
    const filename = issuerId ? `${issuerId}-${ts}` : ts;
    try {
      // The download service types `content` as string and has no `der` entry in its
      // extension/MIME map. Both are fine at runtime as File() accepts a Blob and unknown
      // extensions fall back to text/plain, which is what <DownloadButton> did untyped.
      this.download.miscExtension(filename, content as string, format as keyof Extensions);
      this.flashMessages.info(`Downloading ${filename}`);
    } catch (error) {
      this.flashMessages.danger(errorMessage(error, 'There was a problem downloading. Please try again.'));
    }
    close();
  }
}
