/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from './application';
import { camelize } from '@ember/string';

export default class TotpKeySerializer extends ApplicationSerializer {
  normalizeItems(payload, requestType) {
    if (
      requestType !== 'queryRecord' &&
      payload.data &&
      payload.data.keys &&
      Array.isArray(payload.data.keys)
    ) {
      // if we have data.keys, it's a list of ids, so we map over that
      // and create objects with id's
      return payload.data.keys.map((secret) => ({
        id: secret,
        backend: payload.backend,
      }));
    }

    Object.assign(payload, payload.data);
    delete payload.data;
    return payload;
  }

  serialize(snapshot) {
    // remove all fields that are not relevant to specified key provider
    const { keyFormFields } = snapshot.adapterOptions;
    const json = super.serialize(...arguments);
    Object.keys(json).forEach((key) => {
      if (!keyFormFields.includes(camelize(key))) {
        delete json[key];
      }
    });

    // remove name as it isn't a parameter - it is a part of the request url
    delete json.name;
    return json;
  }
}
