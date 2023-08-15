/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Store from '@ember-data/store';
import { AdapterRegistry } from 'ember-data/adapter';

export interface LdapLibraryAccountStatus {
  account: string;
  available: boolean;
  library: string;
  borrower_client_token?: string;
  borrower_entity_id?: string;
}

export interface LdapLibraryCheckOutCredentials {
  account: string;
  password: string;
  lease_id: string;
  lease_duration: number;
  renewable: boolean;
}

export default interface LdapLibraryAdapter extends AdapterRegistry {
  fetchCheckOutStatus(backend: string, name: string): Promise<Array<LdapLibraryAccountStatus>>;
  checkOutAccount(backend: string, name: string, ttl?: string): Promise<LdapLibraryCheckOutCredentials>;
  checkInAccount(backend: string, name: string, service_account_names: Array<string>): Promise<void>;
}
