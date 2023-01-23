import { underscore } from '@ember/string';
import { keyParamsByType } from 'pki/utils/action-params';
import ApplicationSerializer from '../application';

export default class PkiActionSerializer extends ApplicationSerializer {
  attrs = {
    customTtl: { serialize: false },
    type: { serialize: false },
  };

  serialize(snapshot, requestType) {
    const data = super.serialize(snapshot);
    // requestType is a custom value specified from the pki/action adapter
    const allowedPayloadAttributes = this._allowedParamsByType(requestType, snapshot.record.type);
    if (!allowedPayloadAttributes) return data;

    const payload = {};
    allowedPayloadAttributes.forEach((key) => {
      if ('undefined' !== typeof data[key]) {
        payload[key] = data[key];
      }
    });
    return payload;
  }

  _allowedParamsByType(actionType, type) {
    const keyFields = keyParamsByType(type).map((attrName) => underscore(attrName).toLowerCase());
    const generateProps = [
      'alt_names',
      'common_name',
      'country',
      'exclude_cn_from_sans',
      'format',
      'ip_sans',
      'issuer_name',
      'locality',
      'not_after',
      'not_before_duration',
      'organization',
      'other_sans',
      'ou',
      'permitted_dns_domains',
      'postal_code',
      'private_key_format',
      'province',
      'serial_number',
      'street_address',
      'ttl',
      'type',
      ...keyFields,
    ];
    switch (actionType) {
      case 'import':
        return ['pem_bundle'];
      case 'generate-root':
        return [...generateProps, 'max_path_length'];
      case 'generate-csr':
        return generateProps;
      default:
        // if type doesn't match, serialize all
        return null;
    }
  }
}
