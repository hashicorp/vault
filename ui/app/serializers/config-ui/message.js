/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from '../application';

export default class MessageSerializer extends ApplicationSerializer {
  primaryKey = 'id';

  serialize() {
    const json = super.serialize(...arguments);
    // json.start_time = new Date(json.start_time).toISOString();
    // json.end_time = new Date(json.end_time).toISOString();
    json.link = {
      title: json.link_title,
      href: json.link_href,
    };

    delete json.link_title;
    delete json.link_href;

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
