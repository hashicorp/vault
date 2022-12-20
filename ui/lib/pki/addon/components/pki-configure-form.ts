import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

/**
 * @module PkiConfigureForm
 * PkiConfigureForm components are used to configure a PKI engine mount.
 *
 */

export default class PkiConfigureForm extends Component {
  @tracked configType = '';

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
}
