import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
// TYPES
import Store from '@ember-data/store';
import Router from '@ember/routing/router';
import FlashMessageService from 'vault/services/flash-messages';
import PkiConfigModel from 'vault/models/pki/config';
import { tracked } from '@glimmer/tracking';

interface Args {
  config: PkiConfigModel;
}

/**
 * @module PkiConfigureForm
 * PkiConfigureForm component is used to configure a PKI engine mount.
 * The component shows three options for configuration and handles
 * the save and cancel actions. The sub-forms rendered handle which
 * attributes of the form is shown, based on the formType
 */
export default class PkiConfigureForm extends Component<Args> {
  @service declare readonly store: Store;
  @service declare readonly router: Router;
  @service declare readonly flashMessages: FlashMessageService;
  @tracked formType = '';

  get configTypes() {
    return [
      {
        key: 'import',
        icon: 'download',
        label: 'Import a CA',
        description:
          'Import CA information via a PEM file containing the CA certificate and any private keys, concatenated together, in any order.',
      },
      {
        key: 'generate-root',
        icon: 'file-plus',
        label: 'Generate root',
        description:
          'Generates a new self-signed CA certificate and private key. This generated root will sign its own CRL.',
      },
      {
        key: 'generate-csr',
        icon: 'files',
        label: 'Generate intermediate CSR',
        description:
          'Generate a new CSR for signing, optionally generating a new private key. No new issuer is created by this call.',
      },
    ];
  }

  getFlashMessage(type: string, successful: boolean): string {
    if (type === 'import') {
      return successful
        ? 'Successfully imported the certificate.'
        : 'Could not import the given certificate.';
    }
    // TODO: Fill in messages based on type
    return successful ? 'Configuration successful.' : 'Could not complete configuration';
  }

  shouldUseIssuerEndpoint() {
    const { config } = this.args;
    // To determine which endpoint the config adapter should use,
    // we want to check capabilities on the newer endpoints (those
    // prefixed with "issuers") and use the old path as fallback
    // if user does not have permissions.
    switch (this.formType) {
      case 'import':
        return config.canImportBundle;
      case 'generate-root':
        return config.canGenerateIssuerRoot;
      case 'generate-csr':
        return config.canGenerateIssuerIntermediate;
      default:
        return false;
    }
  }
}
