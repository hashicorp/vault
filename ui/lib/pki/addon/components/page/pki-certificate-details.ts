/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';
import FlashMessages from 'vault/services/flash-messages';
import DownloadService from 'vault/services/download';
import PkiCertificateBaseModel from 'vault/models/pki/certificate/base';

interface Args {
  model: PkiCertificateBaseModel;
  onRevoke?: CallableFunction;
  onBack?: CallableFunction;
}

export default class PkiCertificateDetailsComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessages;
  @service declare readonly download: DownloadService;

  @action
  downloadCert() {
    try {
      const formattedSerial = this.args.model.serialNumber?.replace(/(\s|:)+/g, '-');
      this.download.pem(formattedSerial, this.args.model.certificate);
      this.flashMessages.info('Your download has started.');
    } catch (err) {
      this.flashMessages.danger(errorMessage(err, 'Unable to prepare certificate for download.'));
    }
  }

  @task
  @waitFor
  *revoke() {
    try {
      // the adapter updateRecord method calls the revoke endpoint since it is the only way to update a cert
      yield this.args.model.save();
      this.flashMessages.success('The certificate has been revoked.');
      if (this.args.onRevoke) {
        this.args.onRevoke();
      }
    } catch (error) {
      this.flashMessages.danger(
        errorMessage(error, 'Could not revoke certificate. See Vault logs for details.')
      );
    }
  }
}
