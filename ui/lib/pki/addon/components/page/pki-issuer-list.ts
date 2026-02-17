/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { next } from '@ember/runloop';
import Component from '@glimmer/component';
import { PKI_DEFAULT_EMPTY_STATE_MSG } from 'pki/routes/overview';

import type { PkiReadIssuerResponse } from '@hashicorp/vault-client-typescript';
import type { ParsedCertificateData } from 'vault/utils/parse-pki-cert';

interface BasicDropdown {
  actions: {
    close: CallableFunction;
  };
}
type Issuer = PkiReadIssuerResponse & {
  id: string;
  is_default: boolean;
  serial_number: string;
  isRoot: boolean;
  parsedCertificate: ParsedCertificateData;
};
interface Args {
  issuers: Issuer[];
  mountPoint: string;
  backend: string;
}

export default class PkiIssuerList extends Component<Args> {
  notConfiguredMessage = PKI_DEFAULT_EMPTY_STATE_MSG;

  // To prevent production build bug of passing D.actions to on "click": https://github.com/hashicorp/vault/pull/16983
  @action onLinkClick(D: BasicDropdown) {
    next(() => D.actions.close());
  }
}
