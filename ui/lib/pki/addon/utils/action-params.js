/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * keyParamsByType
 * @param {string} type - refers to `type` attribute on the pki/action model. Should be one of 'exported', 'internal', 'existing', 'kms'
 * @returns array of valid key-related attribute names (camelCase). NOTE: Key params are not used on all action endpoints
 */
export function keyParamsByType(type) {
  let fields = ['key_name', 'key_type', 'key_bits'];
  if (type === 'existing') {
    fields = ['key_ref'];
  } else if (type === 'kms') {
    fields = ['key_name', 'managed_key_name', 'managed_key_id'];
  } else if (type === 'exported') {
    fields = [...fields, 'private_key_format'];
  }
  return fields;
}
