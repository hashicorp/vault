/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */
import type { WithFormFieldsAndValidationsModel } from 'vault/app-types';

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
}
