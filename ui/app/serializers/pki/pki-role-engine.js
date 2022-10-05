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

const EXT_KEY_USAGE = [
  'any',
  'server_auth',
  'client_auth',
  'codes_signing',
  'email_protection',
  'ipsec_end_system',
  'ipsec_tunnel',
  'time_stamping',
  'ocsp_signing',
  'ipsec_user',
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
    const filteredKeyUsage = jsonAsArray.filter(([key]) => {
      return KEY_USAGE.includes(key);
    });
    const filteredExtKeyUsage = jsonAsArray.filter(([key]) => {
      return EXT_KEY_USAGE.includes(key);
    });

    const filteredFlattenedKeyUsage = filteredKeyUsage.flat();
    const filteredFlattenedExtKeyUsage = filteredExtKeyUsage.flat();

    json.key_usage = filteredFlattenedKeyUsage.filter(this.stringsOnly);
    json.ext_key_usage = filteredFlattenedExtKeyUsage.filter(this.stringsOnly);

    filteredFlattenedKeyUsage.filter(this.stringsOnly).forEach((param) => {
      delete json[param];
    });
    filteredFlattenedExtKeyUsage.filter(this.stringsOnly).forEach((param) => {
      delete json[param];
    });

    return json;
  }
}
