import { action } from '@ember/object';
import Component from '@glimmer/component';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';
import { waitFor } from '@ember/test-waiters';

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

export default class PkiCaCertificateImport extends Component {
  @service flashMessages;

  @tracked errorBanner;
  @tracked invalidFormAlert;

  @task
  @waitFor
  *submitForm(event) {
    event.preventDefault();
    try {
      yield this.args.model.save({ adapterOptions: { import: true } });
      this.flashMessages.success('Successfully imported certificate.');
      this.args.onSave();
    } catch (error) {
      this.errorBanner = errorMessage(error);
      this.invalidFormAlert = 'There was a problem importing issuer.';
    }
  }

  @action
  onFileUploaded({ value }) {
    this.args.model.pemBundle = value;
  }

  @action
  cancel() {
    this.args.model.unloadRecord();
    this.args.onCancel();
  }
}
