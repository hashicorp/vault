/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';
import { toLabel } from 'core/helpers/to-label';
import { parseCertificate } from 'vault/utils/parse-pki-cert';

import type FlashMessageService from 'vault/services/flash-messages';
import type DownloadService from 'vault/services/download';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type { ParsedCertificateData } from 'vault/utils/parse-pki-cert';
import type Owner from '@ember/owner';

type CertificateDetails = {
  certificate?: string;
  common_name?: string;
  revocation_time?: number;
  serial_number?: string;
  ca_chain?: string[];
  issuing_ca?: string;
  private_key?: string;
  private_key_type?: string;
};

interface Args {
  certData: CertificateDetails;
  canRevoke: boolean;
  onRevoke?: CallableFunction;
  onBack?: CallableFunction;
}

export default class PkiCertificateDetailsComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly download: DownloadService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  parsedCertificate: ParsedCertificateData;
  @tracked didRevoke = false;

  constructor(owner: Owner, args: Args) {
    super(owner, args);
    this.parsedCertificate = parseCertificate(this.args.certData.certificate || '');
  }

  get displayFields() {
    const fields = [
      'certificate',
      'common_name',
      'serial_number',
      'ca_chain',
      'issuing_ca',
      'private_key',
      'private_key_type',
    ];
    // insert revocation_time after common_name if revoked
    if (this.args.certData.revocation_time || this.didRevoke) {
      fields.splice(2, 0, 'revocation_time');
    }
    return fields;
  }

  isCertificate = (field: string) => ['certificate', 'issuing_ca', 'ca_chain', 'private_key'].includes(field);

  label = (field: string) => {
    const label = toLabel([field]);
    return (
      {
        ca_chain: 'CA chain',
        issuing_ca: 'Issuing CA',
      }[field] || label
    );
  };

  @action
  downloadCert() {
    try {
      const { certificate, serial_number } = this.args.certData;
      const formattedSerial = serial_number?.replace(/(\s|:)+/g, '-') || '';
      this.download.pem(formattedSerial, certificate as string);
      this.flashMessages.info('Your download has started.');
    } catch (err) {
      this.flashMessages.danger(errorMessage(err, 'Unable to prepare certificate for download.'));
    }
  }

  revoke = task(
    waitFor(async () => {
      try {
        const { certificate, serial_number } = this.args.certData;
        // either serial_number or certificate must be provided to revoke but not both
        const payload = serial_number ? { serial_number } : { certificate };
        const { revocation_time } = await this.api.secrets.pkiRevoke(
          this.secretMountPath.currentPath,
          payload
        );
        this.args.certData.revocation_time = revocation_time;
        this.didRevoke = true; // triggers a re-render by adding revocation_time to displayFields
        this.flashMessages.success('The certificate has been revoked.');
        if (this.args.onRevoke) {
          this.args.onRevoke();
        }
      } catch (error) {
        const { message } = await this.api.parseError(
          error,
          'Could not revoke certificate. See Vault logs for details.'
        );
        this.flashMessages.danger(message);
      }
    })
  );
}
