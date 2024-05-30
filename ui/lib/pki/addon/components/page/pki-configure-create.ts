/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import type Store from '@ember-data/store';
import type Router from '@ember/routing/router';
import type FlashMessageService from 'vault/services/flash-messages';
import type PkiActionModel from 'vault/models/pki/action';
import type { Breadcrumb } from 'vault/vault/app-types';

interface Args {
  config: PkiActionModel;
  onCancel: CallableFunction;
  breadcrumbs: Breadcrumb;
}

/**
 * @module PkiConfigureCreate
 * Page::PkiConfigureCreate component is used to configure a PKI engine mount.
 * The component shows three options for configuration and which form
 * is shown. The sub-forms rendered handle rendering the form itself
 * and form submission and cancel actions.
 */
export default class PkiConfigureCreate extends Component<Args> {
  @service declare readonly store: Store;
  @service declare readonly router: Router;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked title = 'Configure PKI';

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
