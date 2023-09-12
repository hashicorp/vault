import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';

import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type AuthService from 'vault/services/auth';
import type LdapLibraryModel from 'vault/models/ldap/library';
import type { LdapLibraryAccountStatus } from 'vault/adapters/ldap/library';

interface Args {
  libraries: Array<LdapLibraryModel>;
  statuses: Array<LdapLibraryAccountStatus>;
  showLibraryColumn: boolean;
}

export default class LdapAccountsCheckedOutComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly router: RouterService;
  @service declare readonly auth: AuthService;

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
    // if disable_check_in_enforcement is true on the library set then all checked out accounts are displayed
    return this.args.statuses.filter((status) => {
      const authEntityId = this.auth.authData?.entity_id;
      const isRoot = !status.borrower_entity_id && !authEntityId; // root user will not have an entity id and it won't be populated on status
      const isEntity = status.borrower_entity_id === authEntityId;
      const library = this.findLibrary(status.library);
      const enforcementDisabled = library.disable_check_in_enforcement === 'Disabled';

      return !status.available && (enforcementDisabled || isEntity || isRoot);
    });
  }

  disableCheckIn = (name: string) => {
    return !this.findLibrary(name).canCheckIn;
  };

  findLibrary(name: string): LdapLibraryModel {
    return this.args.libraries.find((library) => library.name === name) as LdapLibraryModel;
  }

  @task
  @waitFor
  *checkIn() {
    const { library, account } = this.selectedStatus as LdapLibraryAccountStatus;
    try {
      const libraryModel = this.findLibrary(library);
      yield libraryModel.checkInAccount(account);
      this.flashMessages.success(`Successfully checked in the account ${account}.`);
      // transitioning to the current route to trigger the model hook so we can fetch the updated status
      this.router.transitionTo('vault.cluster.secrets.backend.ldap.libraries.library.details.accounts');
    } catch (error) {
      this.selectedStatus = undefined;
      this.flashMessages.danger(`Error checking in the account ${account}. \n ${errorMessage(error)}`);
    }
  }
}
