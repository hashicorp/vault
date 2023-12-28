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
    start_time: { serialize: false },
    end_time: { serialize: false },
  };

  getISODateFormat(date, jsonTime) {
    if (typeof date === 'object') {
      return jsonTime;
    }

    if (typeof date === 'string' && !date.includes('Z')) {
      return new Date(date).toISOString();
    }

    return date;
  }

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
    // it gets to the serialize function. we need to check to see if it's in iso format. if it's not
    // we need to convert it to an ISOString. when the date from the snapshot is an object, it's in a date
    // object format and will need to converted to ISOString. When the date is a string and is not in ISO
    // format, it will need to be converted to an ISOString.
    json.start_time = this.getISODateFormat(snapshot.record.startTime, json.start_time);
    json.end_time = this.getISODateFormat(snapshot.record.endTime, json.end_time);
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
