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

interface File {
  value: string;
  fileName?: string;
  enterAsText: boolean;
}

/**
 * Pki Config Import Component creates a PKI Config record on mount, and cleans it up if dirty on unmount
 */
export default class PkiConfigImportComponent extends Component<Record<string, never>> {
  @service declare readonly store: Store;
  @service declare readonly router: Router;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked configModel: PkiConfigImportModel | null = null;

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
  onFileUploaded(file: File) {
    if (!this.configModel) return;
    this.configModel.pemBundle = file.value;
  }

  @action submitForm(evt: HTMLElementEvent<HTMLFormElement>) {
    evt.preventDefault();
    if (!this.configModel) return;
    this.configModel
      .save()
      .then(() => {
        this.flashMessages.success('Successfully imported the certificate.');
        this.router.transitionTo('vault.cluster.secrets.backend.pki.issuers.index');
      })
      .catch((e) => {
        this.flashMessages.danger(errorMessage(e, 'Could not import the given certificate.'));
      });
  }
}
