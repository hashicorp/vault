/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { decodeString, encodeString } from 'core/utils/b64';
import ApplicationSerializer from '../application';

export default class MessageSerializer extends ApplicationSerializer {
  attrs = {
    link: { serialize: false },
    active: { serialize: false },
  };

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (requestType === 'queryRecord') {
      const transformed = {
        ...payload.data,
        message: decodeString(payload.data.message),
        link_title: payload.data.link.title,
        link_href: payload.data.link.href,
      };
      delete transformed.link;
      return super.normalizeResponse(store, primaryModelClass, transformed, id, requestType);
    }
    return super.normalizeResponse(store, primaryModelClass, payload, id, requestType);
  }

  serialize(snapshot) {
    const json = super.serialize(...arguments);
    json.message = encodeString(json.message);
    json.link = {
      title: json?.link_title || '',
      href: json?.link_href || '',
    };
    // using the snapshot startTime and endTime since the json start and end times are null when
    // it gets to the serialize function.
    json.start_time = snapshot.record?.startTime.includes('Z')
      ? snapshot.record.startTime
      : new Date(snapshot.record.startTime).toISOString();
    json.end_time = snapshot.record?.startTime.includes('Z')
      ? snapshot.record.endTime
      : new Date(snapshot.record.endTime).toISOString();
    delete json?.link_title;
    delete json?.link_href;
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
