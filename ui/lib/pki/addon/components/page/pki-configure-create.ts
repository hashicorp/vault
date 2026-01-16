/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

import type RouterService from '@ember/routing/router';
import type FlashMessageService from 'vault/services/flash-messages';
import type { CapabilitiesMap } from 'vault/vault/app-types';

interface Args {
  capabilities: CapabilitiesMap;
  onCancel: CallableFunction;
}

/**
 * @module PkiConfigureCreate
 * Page::PkiConfigureCreate component is used to configure a PKI engine mount.
 * The component shows three options for configuration and which form
 * is shown. The sub-forms rendered handle rendering the form itself
 * and form submission and cancel actions.
 */
export default class PkiConfigureCreate extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;

  @tracked showActionTypes = true;
  @tracked actionType = '';

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

  @action
  onSave() {
    this.showActionTypes = false;
  }
}
