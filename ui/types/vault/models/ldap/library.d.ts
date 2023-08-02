/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */
import type { WithFormFieldsAndValidationsModel } from 'vault/app-types';
import type { FormField } from 'vault/app-types';
import CapabilitiesModel from '../capabilities';

export default interface LdapRoleModel extends WithFormFieldsAndValidationsModel {
  backend: string;
  name: string;
  service_account_names: string;
  default_ttl: number;
  max_ttl: number;
  disable_check_in_enforcement: string;
}
