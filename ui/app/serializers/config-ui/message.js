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

  getISODateFormat(snapshotDateTime, jsonDateTime) {
    if (typeof snapshotDateTime === 'object') {
      return jsonDateTime;
    }

    // if the snapshot date is in local date time format ("yyyy-MM-dd'T'HH:mm"), we want to ensure
    // it gets converted to an ISOString
    if (typeof snapshotDateTime === 'string' && !snapshotDateTime.includes('Z')) {
      return new Date(snapshotDateTime).toISOString();
    }

    return snapshotDateTime;
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
    // When editing a message with pre-populated dates, this returns a date object. In this case, we would want to use
    // the json date from the serializer. When selecting a date from the datetime-local input date picker, the dates gets
    // set as a date time local string in the model - we would want to convert this local string to an ISOString. Lastly,
    // if this date is not an object and isn’t a local date string, then return the snapshot date, which is set by default
    // values defined on the model.
    json.start_time = this.getISODateFormat(snapshot.record.startTime, json.start_time);
    json.end_time = snapshot.record.endTime
      ? this.getISODateFormat(snapshot.record.endTime, json.end_time)
      : null;
    delete json?.link_title;
    delete json?.link_href;
    return json;
  }

  mapPayload(payload) {
    if (payload.data) {
      if (payload.data?.keys && Array.isArray(payload.data.keys)) {
        return payload.data.keys.map((key) => {
          const data = {
            id: key,
            linkTitle: payload.data.key_info.link?.title,
            linkHref: payload.data.key_info.link?.href,
            ...payload.data.key_info[key],
          };
          if (data.message) data.message = decodeString(data.message);
          return data;
        });
      }
      Object.assign(payload, payload.data);
      delete payload.data;
    }
    return payload;
  }

  extractLazyPaginatedData(payload) {
    return this.mapPayload(payload);
  }
}
