/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { capitalize } from '@ember/string';
import type { KmipWriteRoleRequest } from '@hashicorp/vault-client-typescript';

export default function label(field: keyof KmipWriteRoleRequest) {
  return field.replace('operation_', '').split('_').map(capitalize).join(' ');
}
