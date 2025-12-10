/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

import type FlashMessageService from 'vault/services/flash-messages';
import type AuthService from 'vault/services/auth';
import type { LdapLibrary } from 'vault/secrets/ldap';
import type { LdapLibraryAccountStatus } from 'vault/adapters/ldap/library';
import type { CapabilitiesMap } from 'vault/app-types';
import type CapabilitiesService from 'vault/services/capabilities';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';

interface Args {
  libraries: Array<LdapLibrary>;
  capabilities: CapabilitiesMap;
  statuses: Array<LdapLibraryAccountStatus>;
  showLibraryColumn: boolean;
  onCheckInSuccess: CallableFunction;
  isLoadingStatuses?: boolean;
}

export default class LdapAccountsCheckedOutComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly auth: AuthService;
  @service declare readonly capabilities: CapabilitiesService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  @tracked selectedStatus: LdapLibraryAccountStatus | undefined;

  get columns() {
    const columns = [{ label: 'Account' }, { label: 'Action' }];
    if (this.args.showLibraryColumn) {
      columns.splice(1, 0, { label: 'Library' });
    }
    return columns;
  }

  get filteredAccounts() {
    // filter status to only show checked out accounts associated to the current user
    // if disable_check_in_enforcement is true on the library then all checked out accounts are displayed
    return this.args.statuses.filter((status) => {
      const authEntityId = this.auth.authData?.entityId;
      const isRoot = !status.borrower_entity_id && !authEntityId; // root user will not have an entity id and it won't be populated on status
      const isEntity = status.borrower_entity_id === authEntityId;
      const library = this.findLibrary(status.library);

      return !status.available && (library.disable_check_in_enforcement || isEntity || isRoot);
    });
  }

  disableCheckIn = (libraryName: string) => {
    const { completeLibraryName: name } = this.findLibrary(libraryName);
    const { currentPath: backend } = this.secretMountPath;
    const path = this.capabilities.pathFor('ldapLibraryCheckIn', { backend, name });
    const { canUpdate } = this.args.capabilities[path] || {};
    return !canUpdate;
  };

  findLibrary(name: string): LdapLibrary {
    return this.args.libraries.find((library) => library.completeLibraryName === name) as LdapLibrary;
  }

  checkIn = task(
    waitFor(async () => {
      const { library, account } = this.selectedStatus as LdapLibraryAccountStatus;
      try {
        const { completeLibraryName } = this.findLibrary(library);
        const payload = { service_account_names: [account] };
        await this.api.secrets.ldapLibraryForceCheckIn(
          completeLibraryName,
          this.secretMountPath.currentPath,
          payload
        );
        this.flashMessages.success(`Successfully checked in the account ${account}.`);
        this.args.onCheckInSuccess();
      } catch (error) {
        const { message } = await this.api.parseError(error);
        this.selectedStatus = undefined;
        this.flashMessages.danger(`Error checking in the account ${account}. \n ${message}`);
      }
    })
  );
}
