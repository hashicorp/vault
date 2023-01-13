import ApplicationSerializer from '../application';

export default class PkiActionSerializer extends ApplicationSerializer {
  attrs = {
    customTtl: { serialize: false },
    type: { serialize: false },
  };

  serialize(snapshot, requestType) {
    const data = super.serialize(snapshot);
    // requestType is a custom value specified from the pki/action adapter
    const allowedPayloadAttributes = this._allowedParamsByType(requestType);
    if (!allowedPayloadAttributes) return data;

    const payload = {};
    allowedPayloadAttributes.forEach((key) => {
      if ('undefined' !== typeof data[key]) {
        payload[key] = data[key];
      }
    });
    return payload;
  }

  _allowedParamsByType(formType) {
    switch (formType) {
      case 'import':
        return ['pem_bundle'];
      case 'generate-root':
        return [
          'alt_names',
          'common_name',
          'country',
          'exclude_cn_from_sans',
          'format',
          'ip_sans',
          'issuer_name',
          'key_bits',
          'key_name',
          'key_ref',
          'key_type',
          'locality',
          'managed_key_id',
          'managed_key_name',
          'max_path_length',
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
          'type',
        ];
      default:
        // if type doesn't match, serialize all
        return null;
    }
  }
}
