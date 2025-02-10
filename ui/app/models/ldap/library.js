/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import { withModelValidations } from 'vault/decorators/model-validations';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

const validations = {
  name: [{ type: 'presence', message: 'Library name is required.' }],
  service_account_names: [{ type: 'presence', message: 'At least one service account is required.' }],
};
const formFields = ['name', 'service_account_names', 'ttl', 'max_ttl', 'disable_check_in_enforcement'];

@withModelValidations(validations)
@withFormFields(formFields)
export default class LdapLibraryModel extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord
  @attr('string') path_to_library; // ancestral path to the library added in the adapter (only exists for nested libraries)

  @attr('string', {
    label: 'Library name',
    editDisabled: true,
  })
  name;

  @attr('string', {
    editType: 'stringArray',
    label: 'Accounts',
    subText:
      'The names of all the accounts that can be checked out from this set. These accounts must only be used by Vault, and may only be in one set.',
  })
  service_account_names;

  @attr({
    editType: 'ttl',
    label: 'Default lease TTL',
    detailsLabel: 'TTL',
    helperTextDisabled: 'Vault will use the default lease duration.',
    defaultValue: '24h',
    defaultShown: 'Engine default',
  })
  ttl;

  @attr({
    editType: 'ttl',
    label: 'Max lease TTL',
    detailsLabel: 'Max TTL',
    helperTextDisabled: 'Vault will use the default lease duration.',
    defaultValue: '24h',
    defaultShown: 'Engine default',
  })
  max_ttl;

  // this is a boolean from the server but is transformed in the serializer to display as Disabled or Enabled
  @attr('string', {
    editType: 'radio',
    label: 'Check-in enforcement',
    subText:
      'When enabled, accounts must be checked in by the entity or client token that checked them out. If disabled, anyone with the right permission can check the account back in.',
    possibleValues: ['Disabled', 'Enabled'],
    defaultValue: 'Enabled',
  })
  disable_check_in_enforcement;

  get completeLibraryName() {
    // if there is a path_to_library then the name is hierarchical
    // and we must concat the ancestors with the leaf name to get the full library path
    return this.path_to_library ? this.path_to_library + this.name : this.name;
  }

  get displayFields() {
    return this.formFields.filter((field) => field.name !== 'service_account_names');
  }

  @lazyCapabilities(apiPath`${'backend'}/library/${'name'}`, 'backend', 'name') libraryPath;
  @lazyCapabilities(apiPath`${'backend'}/library/${'name'}/status`, 'backend', 'name') statusPath;
  @lazyCapabilities(apiPath`${'backend'}/library/${'name'}/check-out`, 'backend', 'name') checkOutPath;
  @lazyCapabilities(apiPath`${'backend'}/library/${'name'}/check-in`, 'backend', 'name') checkInPath;

  get canCreate() {
    return this.libraryPath.get('canCreate') !== false;
  }
  get canDelete() {
    return this.libraryPath.get('canDelete') !== false;
  }
  get canEdit() {
    return this.libraryPath.get('canUpdate') !== false;
  }
  get canRead() {
    return this.libraryPath.get('canRead') !== false;
  }
  get canList() {
    return this.libraryPath.get('canList') !== false;
  }
  get canReadStatus() {
    return this.statusPath.get('canRead') !== false;
  }
  get canCheckOut() {
    return this.checkOutPath.get('canUpdate') !== false;
  }
  get canCheckIn() {
    return this.checkInPath.get('canUpdate') !== false;
  }

  fetchStatus() {
    return this.store.adapterFor('ldap/library').fetchStatus(this.backend, this.name);
  }
  checkOutAccount(ttl) {
    return this.store.adapterFor('ldap/library').checkOutAccount(this.backend, this.name, ttl);
  }
  checkInAccount(account) {
    return this.store.adapterFor('ldap/library').checkInAccount(this.backend, this.name, [account]);
  }
}
