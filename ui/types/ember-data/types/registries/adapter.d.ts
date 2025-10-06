/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Application from 'vault/adapters/application';
import Adapter from 'ember-data/adapter';
import ModelRegistry from 'ember-data/types/registries/model';

import ClientsActivityAdapter from 'vault/vault/adapters/clients/activity';
import LdapLibraryAdapter from 'vault/adapters/ldap/library';
import LdapRoleAdapter from 'vault/adapters/ldap/role';
import PkiIssuerAdapter from 'vault/adapters/pki/issuer';
import PkiTidyAdapter from 'vault/adapters/pki/tidy';
import SyncAssociationAdapter from 'vault/adapters/sync/association';
import SyncDestinationAdapter from 'vault/adapters/sync/destination';

/**
 * Catch-all for ember-data.
 */
export default interface AdapterRegistry {
  'clients/activity': ClientsActivityAdapter;
  'ldap/library': LdapLibraryAdapter;
  'ldap/role': LdapRoleAdapter;
  'pki/issuer': PkiIssuerAdapter;
  'pki/tidy': PkiTidyAdapter;
  'sync/destination': SyncDestinationAdapter;
  'sync/association': SyncAssociationAdapter;
  application: Application;
  [key: keyof ModelRegistry]: Adapter;
}
