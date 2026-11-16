/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';

import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type { TtlEvent } from 'vault/app-types';
import type { LdapLibrary, LdapLibraryAccountStatus } from 'vault/secrets/ldap';

interface Args {
  library: LdapLibrary;
  statuses: Array<LdapLibraryAccountStatus>;
}

export default class LdapLibraryDetailsAccountsPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly secretMountPath: SecretMountPath;

  @tracked showCheckOutPrompt = false;
  @tracked checkOutTtl: string | null = null;
  @tracked isTtlValid = true;

  get cliCommand() {
    return `vault lease renew ${this.secretMountPath.currentPath}/library/${this.args.library.name}/check-out/:lease_id`;
  }
  // The picker starts over each time the prompt opens, so a value entered earlier must not carry over.
  @action
  toggleCheckOutPrompt(isOpen: boolean) {
    this.showCheckOutPrompt = isOpen;
    this.checkOutTtl = null;
    this.isTtlValid = true;
  }
  @action
  setTtl(data: TtlEvent) {
    this.checkOutTtl = data.timeString;
    this.isTtlValid = data.isValid;
  }
  @action
  checkOut() {
    // the picker explains an over-limit TTL, so it only needs to stop the check-out
    if (!this.isTtlValid) return;
    this.router.transitionTo('vault.cluster.secrets.backend.ldap.libraries.library.check-out', {
      queryParams: { ttl: this.checkOutTtl },
    });
  }
}
