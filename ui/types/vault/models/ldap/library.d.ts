/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */
import type { WithFormFieldsAndValidationsModel } from 'vault/app-types';
import type { FormField } from 'vault/app-types';
import CapabilitiesModel from '../capabilities';
import type {
  LdapLibraryAccountStatus,
  LdapLibraryCheckOutCredentials,
} from 'vault/vault/adapters/ldap/library';

export default interface LdapLibraryModel extends WithFormFieldsAndValidationsModel {
  backend: string;
  name: string;
  service_account_names: string;
  default_ttl: number;
  max_ttl: number;
  disable_check_in_enforcement: string;
  get displayFields(): Array<FormField>;
  libraryPath: CapabilitiesModel;
  statusPath: CapabilitiesModel;
  checkOutPath: CapabilitiesModel;
  checkInPath: CapabilitiesModel;
  get canCreate(): boolean;
  get canDelete(): boolean;
  get canEdit(): boolean;
  get canRead(): boolean;
  get canList(): boolean;
  get canReadStatus(): boolean;
  get canCheckOut(): boolean;
  get canCheckIn(): boolean;
  fetchStatus(): Promise<Array<LdapLibraryAccountStatus>>;
  checkOutAccount(ttl?: string): Promise<LdapLibraryCheckOutCredentials>;
  checkInAccount(account: string): Promise<void>;
}
