import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
// TYPES
import Store from '@ember-data/store';
import Router from '@ember/routing/router';
import FlashMessageService from 'vault/services/flash-messages';
import { action } from '@ember/object';
import PkiActionModel from 'vault/models/pki/action';

interface Args {
  config: PkiActionModel;
}

/**
 * @module PkiConfigureForm
 * PkiConfigureForm component is used to configure a PKI engine mount.
 * The component shows three options for configuration and which form
 * is shown. The sub-forms rendered handle rendering the form itself
 * and form submission and cancel actions.
 */
export default class PkiConfigureForm extends Component<Args> {
  @service declare readonly store: Store;
  @service declare readonly router: Router;
  @service declare readonly flashMessages: FlashMessageService;

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

  @action cancel() {
    this.args.config.rollbackAttributes();
    this.router.transitionTo('vault.cluster.secrets.backend.pki.overview');
  }
}
