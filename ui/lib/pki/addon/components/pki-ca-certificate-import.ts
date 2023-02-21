import { action } from '@ember/object';
import Component from '@glimmer/component';
import FlashMessageService from 'vault/services/flash-messages';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';
import PkiCertificateBaseModel from 'vault/models/pki/certificate/base';
import PkiActionModel from 'vault/models/pki/action';

/**
 * @module PkiCaCertificateImport
 * PkiCaCertificateImport components are used to import PKI CA certificates and keys via pem_bundle.
 * https://github.com/hashicorp/vault/blob/main/website/content/api-docs/secret/pki.mdx#import-ca-certificates-and-keys
 *
 * @example
 * ```js
 * <PkiCaCertificateImport @model={{this.model}} />
 * ```
 *
 * @param {Object} model - certificate model from route
 * @callback onCancel - Callback triggered when cancel button is clicked.
 * @callback onSubmit - Callback triggered on submit success.
 */

interface Args {
  onSave: CallableFunction;
  onCancel: CallableFunction;
  model: PkiCertificateBaseModel | PkiActionModel;
  adapterOptions: object | undefined;
}

export default class PkiCaCertificateImport extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;

  @tracked errorBanner = '';

  @task
  @waitFor
  *submitForm(event: Event) {
    event.preventDefault();
    try {
      yield this.args.model.save({ adapterOptions: this.args.adapterOptions });
      this.flashMessages.success('Successfully imported certificate.');
      this.args.onSave();
    } catch (error) {
      this.errorBanner = errorMessage(error);
    }
  }

  @action
  onFileUploaded({ value }: { value: string }) {
    this.args.model.pemBundle = value;
  }

  @action
  cancel() {
    this.args.model.unloadRecord();
    this.args.onCancel();
  }
}
