/**
 * Copyright IBM Corp. 2016, 2025
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

  get cliCommand() {
    return `vault lease renew ${this.secretMountPath.currentPath}/library/${this.args.library.name}/check-out/:lease_id`;
  }
  @action
  setTtl(data: TtlEvent) {
    this.checkOutTtl = data.timeString;
  }
  @action
  checkOut() {
    this.router.transitionTo('vault.cluster.secrets.backend.ldap.libraries.library.check-out', {
      queryParams: { ttl: this.checkOutTtl },
    });
  }
}
