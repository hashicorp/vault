/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { decodeString, encodeString } from 'core/utils/b64';
import ApplicationSerializer from '../application';

export default class MessageSerializer extends ApplicationSerializer {
  primaryKey = 'id';

  normalizeResponse(store, primaryModelClass, payload) {
    if (payload.data.data && payload.data.data?.message) {
      payload.data.data.message = decodeString(payload.data.data.message);
      payload.data = {
        id: payload.data.id,
        linkTitle: payload.data.data.link?.title,
        linkHref: payload.data.data.link?.href,
        ...payload.data.data,
      };
      delete payload.data.data;
    }
    return super.normalizeResponse(...arguments);
  }

  serialize() {
    const json = super.serialize(...arguments);
    json.message = encodeString(json.message);
    json.link = {
      title: json.link_title,
      href: json.link_href,
    };

    delete json.link_title;
    delete json.link_href;
    delete json.active;

    return json;
  }

  extractLazyPaginatedData(payload) {
    if (payload.data) {
      if (payload.data?.keys && Array.isArray(payload.data.keys)) {
        return payload.data.keys.map((key) => {
          return {
            id: key,
            linkTitle: payload.data.key_info.link?.title,
            linkHref: payload.data.key_info.link?.href,
            ...payload.data.key_info[key],
          };
        });
      }
      Object.assign(payload, payload.data);
      delete payload.data;
    }
    return payload;
  }
}
