/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { underscore } from '@ember/string';
import { keyParamsByType } from 'pki/utils/action-params';
import ApplicationSerializer from '../application';
import { parseCertificate } from 'vault/utils/parse-pki-cert';

export default class PkiActionSerializer extends ApplicationSerializer {
  attrs = {
    customTtl: { serialize: false },
    type: { serialize: false },
  };

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (payload.data.certificate) {
      // Parse certificate back from the API and add to payload
      const parsedCert = parseCertificate(payload.data.certificate);
      const data = {
        ...payload.data,
        common_name: parsedCert.common_name,
        parsed_certificate: parsedCert,
      };
      return super.normalizeResponse(store, primaryModelClass, { ...payload, data }, id, requestType);
    }
    return super.normalizeResponse(...arguments);
  }

  serialize(snapshot, requestType) {
    const data = super.serialize(snapshot);
    // requestType is a custom value specified from the pki/action adapter
    const allowedPayloadAttributes = this._allowedParamsByType(requestType, snapshot.record.type);
    if (!allowedPayloadAttributes) return data;
    // the backend expects the subject's serial number param to be 'serial_number'
    // we label it as subject_serial_number to differentiate from the vault generated UUID
    data.serial_number = data.subject_serial_number;

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
    const commonProps = [
      'alt_names',
      'common_name',
      'country',
      'exclude_cn_from_sans',
      'format',
      'ip_sans',
      'locality',
      'organization',
      'other_sans',
      'ou',
      'postal_code',
      'province',
      'serial_number',
      'street_address',
      'type',
      'uri_sans',
      ...keyFields,
    ];
    switch (actionType) {
      case 'import':
        return ['pem_bundle'];
      case 'generate-root':
        return [
          ...commonProps,
          'issuer_name',
          'max_path_length',
          'not_after',
          'not_before_duration',
          'permitted_dns_domains',
          'private_key_format',
          'ttl',
        ];
      case 'rotate-root':
        return [
          ...commonProps,
          'issuer_name',
          'max_path_length',
          'not_after',
          'not_before_duration',
          'permitted_dns_domains',
          'private_key_format',
          'ttl',
        ];
      case 'generate-csr':
        return [...commonProps, 'add_basic_constraints'];
      case 'sign-intermediate':
        return ['common_name', 'issuer_name', 'csr'];
      default:
        // if type doesn't match, serialize all
        return null;
    }
  }
}
