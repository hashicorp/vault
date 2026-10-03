/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { KmipWriteRoleRequest, KmipWriteRoleRequestToJSONTyped } from '@hashicorp/vault-client-typescript';

export default function operationGroups() {
  // the client helper returns an object with all available properties from the schema
  // we could pass in a role to set the values but in this case we are only interested in the keys
  const role = KmipWriteRoleRequestToJSONTyped({} as KmipWriteRoleRequest);
  const objects = [
    'operation_create',
    'operation_activate',
    'operation_get',
    'operation_locate',
    'operation_rekey',
    'operation_revoke',
    'operation_destroy',
  ];
  const attributes = ['operation_add_attribute', 'operation_get_attributes'];
  const server = ['operation_discover_versions'];
  const notOther = [...objects, ...attributes, ...server, 'operation_all', 'operation_none'];
  const other = Object.keys(role).filter((key) => key.startsWith('operation_') && !notOther.includes(key));

  return {
    'Managed Cryptographic Objects': objects,
    'Object Attributes': attributes,
    Server: server,
    Other: other,
  };
}
