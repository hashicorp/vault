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
  @tracked invalidFormAlert = '';
  @tracked modelValidations = null;

  get verb() {
    return this.args.type === 'sign' ? 'sign' : 'generate';
  }

  @task
  *save(evt: Event) {
    evt.preventDefault();
    this.errorBanner = '';
    const { model, onSuccess } = this.args;
    const { isValid, state, invalidFormMessage } = model.validate();

    this.modelValidations = isValid ? null : state;
    this.invalidFormAlert = invalidFormMessage;

    if (!isValid) return;

    try {
      yield model.save();
      onSuccess();
    } catch (err) {
      this.errorBanner = errorMessage(err, `Could not ${this.verb} certificate. See Vault logs for details.`);
    }
  }

  @action cancel() {
    this.args.model.unloadRecord();
    this.router.transitionTo('vault.cluster.secrets.backend.pki.roles.role.details');
  }
}
