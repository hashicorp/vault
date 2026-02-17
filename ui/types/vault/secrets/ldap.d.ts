/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

export type LdapRole = {
  name: string;
  type: string;
  completeRoleName: string;
};

export type LdapStaticRole = LdapRole & {
  dn: string;
  username: string;
  rotation_period: string;
};

export type LdapDynamicRole = LdapRole & {
  default_ttl: string;
  max_ttl: string;
  username_template: string;
  creation_ldif: string;
  deletion_ldif: string;
  rollback_ldif: string;
};

export interface LdapStaticRoleCredentials {
  dn: string;
  last_vault_rotation: string;
  password: string;
  last_password: string;
  rotation_period: number;
  ttl: number;
  username: string;
  type: string;
}

export interface LdapDynamicRoleCredentials {
  distinguished_names: Array<string>;
  password: string;
  username: string;
  lease_id: string;
  lease_duration: string;
  renewable: boolean;
  type: string;
}

export type LdapLibrary = {
  name: string;
  completeLibraryName: string;
  service_account_names: string[];
  ttl: string;
  max_ttl: string;
  disable_check_in_enforcement: boolean;
};

export type LdapLibraryAccountStatusResponse = Record<
  string,
  { available: boolean; borrower_client_token?: string; borrower_entity_id?: string }
>;

export type LdapLibraryAccountStatus = {
  account: string;
  available: boolean;
  library: string;
  borrower_client_token?: string;
  borrower_entity_id?: string;
};
