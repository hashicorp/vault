import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';

import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type LdapLibraryModel from 'vault/models/ldap/library';
import type { LdapLibraryAccountStatus } from 'vault/vault/adapters/ldap/library';
import { TtlEvent } from 'vault/vault/app-types';

interface Args {
  library: LdapLibraryModel;
  statuses: Array<LdapLibraryAccountStatus>;
}

export default class LdapLibraryDetailsAccountsPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly router: RouterService;

  @tracked showCheckOutPrompt = false;
  @tracked checkOutTtl: string | null = null;

  get cliCommand() {
    return `vault lease renew ad/library/${this.args.library.name}/check-out/:lease_id`;
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
