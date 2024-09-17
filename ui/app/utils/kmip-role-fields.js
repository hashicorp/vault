/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { removeManyFromArray } from 'vault/helpers/remove-from-array';

export const operationFields = (fieldNames) => {
  if (!Array.isArray(fieldNames)) {
    throw new Error('fieldNames must be an array');
  }
  return fieldNames.filter((key) => key.startsWith('operation'));
};

export const operationFieldsWithoutSpecial = (fieldNames) => {
  const opFields = operationFields(fieldNames);
  return removeManyFromArray(opFields, ['operationAll', 'operationNone']);
};

export const nonOperationFields = (fieldNames) => {
  const opFields = operationFields(fieldNames);
  return removeManyFromArray(fieldNames, opFields);
};

export const tlsFields = () => {
  return ['tlsClientKeyBits', 'tlsClientKeyType', 'tlsClientTtl'];
};
