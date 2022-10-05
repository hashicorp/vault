import ApplicationSerializer from '../application';

const KEY_USAGE = [
  'digital_signature',
  'key_agreement',
  'key_encipherment',
  'content_commitment',
  'data_encripherment',
  'cert_sign',
  'crl_sign',
  'encipher_only',
  'decipher_only',
];

export default class PkiRoleEngineSerializer extends ApplicationSerializer {
  stringsOnly(value) {
    if (typeof value === 'string') {
      return value;
    }
  }

  serialize() {
    const json = super.serialize(...arguments);
    const jsonAsArray = Object.entries(json);
    /*
      Handle adding the KEY_USAGE checkboxes and EXT_KEY_USAGE checkboxes to their respect key_usage and ext_key_usage ['list: []'] https://www.vaultproject.io/api-docs/secret/pki#key_usage
      Process:
      1. Turn the object Json into a key/value array of arrays so we can filter over it.
      2. This returns as the const filtered = [['content_commitment',true],['crl_sign',true]]
      3. Flatten this array of nested arrays to return const filteredFlatted = ['content_commitment',true,'crl_sign',true]
      4. Lastly filter on the filteredFlattened to return only strings so that we add the the key_usage param = ['content_commitment','crl_sign']
      5. cleanup: remove from model the unused params via the delete operation. 
    */
    const filtered = jsonAsArray.filter(([key]) => {
      return KEY_USAGE.includes(key);
    });
    const filteredFlattened = filtered.flat();

    json.key_usage = filteredFlattened.filter(this.stringsOnly);

    filteredFlattened.filter(this.stringsOnly).forEach((param) => {
      delete json[param];
    });

    // empty arrays are being removed from serialized json
    // ensure that they are sent to the server, otherwise removing items will not be persisted
    json.auth_method_accessors = json.auth_method_accessors || [];
    json.auth_method_types = json.auth_method_types || [];
    return this.transformHasManyKeys(json, 'server');
  }
}
