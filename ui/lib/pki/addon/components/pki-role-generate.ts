import Component from '@glimmer/component';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { service } from '@ember/service';
import Router from '@ember/routing/router';
import Store from '@ember-data/store';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';
import FlashMessageService from 'vault/services/flash-messages';
import DownloadService from 'vault/services/download';
import PkiCertificateGenerateModel from 'vault/models/pki/certificate/generate';

interface Args {
  onSuccess: CallableFunction;
  model: PkiCertificateGenerateModel;
  type: string;
}

export default class PkiRoleGenerate extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly store: Store;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly download: DownloadService;

  @tracked errorBanner = '';

  transitionToRole() {
    this.router.transitionTo('vault.cluster.secrets.backend.pki.roles.role.details');
  }

  get verb() {
    return this.args.type === 'sign' ? 'sign' : 'generate';
  }

  @task
  *save(evt: Event) {
    evt.preventDefault();
    this.errorBanner = '';
    const { model, onSuccess } = this.args;
    try {
      yield model.save();
      onSuccess();
    } catch (err) {
      this.errorBanner = errorMessage(err, `Could not ${this.verb} certificate. See Vault logs for details.`);
    }
  }

  @task
  *revoke() {
    try {
      yield this.args.model.destroyRecord();
      this.flashMessages.success('The certificate has been revoked.');
      this.transitionToRole();
    } catch (err) {
      this.errorBanner = errorMessage(err, 'Could not revoke certificate. See Vault logs for details.');
    }
  }

  @action downloadCert() {
    try {
      const formattedSerial = this.args.model.serialNumber?.replace(/(\s|:)+/g, '-');
      this.download.pem(formattedSerial, this.args.model.certificate);
      this.flashMessages.info('Your download has started.');
    } catch (err) {
      this.flashMessages.danger(errorMessage(err, 'Unable to prepare certificate for download.'));
    }
  }

  @action cancel() {
    this.args.model.unloadRecord();
    this.transitionToRole();
  }
}
