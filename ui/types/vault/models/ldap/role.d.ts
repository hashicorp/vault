/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */
import type { WithFormFieldsAndValidationsModel } from 'vault/app-types';
import type { FormField } from 'vault/app-types';
import CapabilitiesModel from '../capabilities';
import { LdapDynamicRoleCredentials, LdapStaticRoleCredentials } from 'ldap/routes/roles/role/credentials';
export default interface LdapRoleModel extends WithFormFieldsAndValidationsModel {
  type: string;
  backend: string;
  name: string;
  dn: string;
  username: string;
  rotation_period: string;
  default_ttl: string;
  max_ttl: string;
  username_template: string;
  creation_ldif: string;
  rollback_ldif: string;
  get isStatic(): string;
  get isDynamic(): string;
  get fieldsForType(): Array<string>;
  get displayFields(): Array<FormField>;
  get roleUri(): string;
  get credsUri(): string;
  rolePath: CapabilitiesModel;
  credsPath: CapabilitiesModel;
  staticRotateCredsPath: CapabilitiesModel;
  get canCreate(): boolean;
  get canDelete(): boolean;
  get canEdit(): boolean;
  get canRead(): boolean;
  get canList(): boolean;
  get canReadCreds(): boolean;
  get canRotateStaticCreds(): boolean;
  fetchCredentials(): Promise<LdapDynamicRoleCredentials | LdapStaticRoleCredentials>;
  rotateStaticPassword(): Promise<void>;
}
