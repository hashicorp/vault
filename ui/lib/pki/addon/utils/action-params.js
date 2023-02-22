/**
 * keyParamsByType
 * @param {string} type - refers to `type` attribute on the pki/action model. Should be one of 'exported', 'internal', 'existing', 'kms'
 * @returns array of valid key-related attribute names (camelCase). NOTE: Key params are not used on all action endpoints
 */
export function keyParamsByType(type) {
  let fields = ['keyName', 'keyType', 'keyBits'];
  if (type === 'existing') {
    fields = ['keyRef'];
  } else if (type === 'kms') {
    fields = ['keyName', 'managedKeyName', 'managedKeyId'];
  } else if (type === 'exported') {
    fields = ['keyName', 'keyType', 'keyBits', 'privateKeyFormat'];
  }
  return fields;
}
