/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { next } from '@ember/runloop';
import Component from '@glimmer/component';
import { PKI_DEFAULT_EMPTY_STATE_MSG } from 'pki/routes/overview';
import type PkiIssuerModel from 'vault/models/pki/issuer';

interface BasicDropdown {
  actions: {
    close: CallableFunction;
  };
}
interface Args {
  issuers: PkiIssuerModel[];
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
