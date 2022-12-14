import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
// Types
import { HTMLElementEvent } from 'forms';
import PkiConfigImportModel from 'vault/models/pki/config/import';
import Store from '@ember-data/store';
import Router from '@ember/routing/router';
import FlashMessageService from 'vault/services/flash-messages';
import errorMessage from 'vault/utils/error-message';

// TODO: validate with Alex what we're looking for here
// https://developer.hashicorp.com/vault/api-docs/secret/pki#parameters-15
// https://polarssl.org/kb/cryptography/asn1-key-structures-in-der-and-pem/
function getCertType(contents: string): string {
  if (contents.startsWith('-----BEGIN CERTIFICATE-----')) {
    return 'certificate';
  }
  return 'pem';
}

interface File {
  value: string;
  fileName?: string;
  enterAsText: boolean;
}

export default class PkiConfigImportComponent extends Component<Record<string, never>> {
  @service declare readonly store: Store;
  @service declare readonly router: Router;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked configModel: PkiConfigImportModel | null = null;
  @tracked file: File = { value: '', enterAsText: false };

  constructor(owner: unknown, args: Record<string, never>) {
    super(owner, args);
    const model = this.store.createRecord('pki/config/import');
    this.configModel = model;
  }

  willDestroy() {
    super.willDestroy();
    const config = this.configModel;
    // error is thrown when you attempt to unload a record that is inFlight (isSaving)
    if ((config?.isNew || config?.hasDirtyAttributes) && !config?.isSaving) {
      config.unloadRecord();
    }
  }

  @action
  onFileUploaded(_: unknown, file: File) {
    const type = getCertType(file.value);
    this.file = file;
    if (!this.configModel) return;
    if (type === 'pem') {
      this.configModel.pemBundle = file.value;
      this.configModel.certificate = undefined;
    } else {
      this.configModel.certificate = file.value;
      this.configModel.pemBundle = undefined;
    }
  }

  @action submitForm(evt: HTMLElementEvent<HTMLFormElement>) {
    evt.preventDefault();
    if (!this.configModel) return;
    this.configModel
      .save()
      .then(() => {
        this.router.transitionTo('vault.cluster.secrets.backend.pki.issuers.index');
      })
      .catch((e) => {
        this.flashMessages.danger(errorMessage(e, 'Could not import the given certificate.'));
      });
  }
}
