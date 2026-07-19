/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { SetupSteps } from '../implementation-select';

import type { AcmeAccountsRouteModel } from 'pki/routes/external/acme-accounts';
import type { PkiExternalCaReadConfigAcmeAccountResponse } from '@hashicorp/vault-client-typescript';

interface Args {
  model: AcmeAccountsRouteModel;
}

export default class ExternalPkiPageAcmeAccountsComponent extends Component<Args> {
  @tracked showFlyout = false;
  @tracked flyoutAcct: PkiExternalCaReadConfigAcmeAccountResponse | undefined;

  acmeConfig = SetupSteps.ACME_CONFIG;
  accountKeysColumns = [
    { key: 'key_creation_date', label: 'Creation date', isSortable: true },
    { key: 'key_type', label: 'Type' },
    { key: 'key_version', label: 'Version', isSortable: true },
  ];

  @action
  openFlyout(config: PkiExternalCaReadConfigAcmeAccountResponse) {
    this.showFlyout = true;
    this.flyoutAcct = config;
  }

  @action
  closeFlyout() {
    this.showFlyout = false;
    this.flyoutAcct = undefined;
  }

  get flyoutActiveKey() {
    return this.flyoutAcct ? this.findActiveKey(this.flyoutAcct) : undefined;
  }

  get inactiveKeysTable() {
    const keys = this.getAccountKeys(this.flyoutAcct);
    const activeVersion = this.flyoutAcct?.active_key_version;
    return Object.fromEntries(keys.filter(([_, key]) => key.key_version !== activeVersion));
  }

  // Template helpers
  getAccountKeys = (acmeAccount?: PkiExternalCaReadConfigAcmeAccountResponse) => {
    return acmeAccount?.account_keys ? Object.entries(acmeAccount.account_keys) : [];
  };

  findActiveKey = (acmeAccount: PkiExternalCaReadConfigAcmeAccountResponse) => {
    const keys = this.getAccountKeys(acmeAccount);
    return keys.find(([_, key]) => key.key_version === acmeAccount.active_key_version)?.[1];
  };

  formatConfigDetails = (acmeAccount: PkiExternalCaReadConfigAcmeAccountResponse) => {
    const activeKey = this.findActiveKey(acmeAccount);
    return { key_type: activeKey?.key_type, ...acmeAccount };
  };

  hasKeyHistory = (acmeAccount: PkiExternalCaReadConfigAcmeAccountResponse): boolean =>
    this.getAccountKeys(acmeAccount).length > 1;

  isObject = (value: unknown): value is Record<string, unknown> =>
    typeof value === 'object' && value !== null;
}
