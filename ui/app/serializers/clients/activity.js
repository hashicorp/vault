/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationSerializer from '../application';
import { formatISO } from 'date-fns';
import { formatByMonths, formatByNamespace, homogenizeClientNaming } from 'core/utils/client-count-utils';
export default class ActivitySerializer extends ApplicationSerializer {
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (payload.id === 'no-data') {
      return super.normalizeResponse(store, primaryModelClass, payload, id, requestType);
    }
    const response_timestamp = formatISO(new Date());
    const transformedPayload = {
      ...payload,
      response_timestamp,
      by_namespace: formatByNamespace(payload.data.by_namespace),
      by_month: formatByMonths(payload.data.months),
      total: homogenizeClientNaming(payload.data.total),
    };
    delete payload.data.by_namespace;
    delete payload.data.months;
    delete payload.data.total;
    return super.normalizeResponse(store, primaryModelClass, transformedPayload, id, requestType);
  }
}
/* 
SAMPLE PAYLOAD BEFORE/AFTER:

payload.data.by_namespace = [
  {
    namespace_id: '5SWT8',
    namespace_path: 'namespacelonglonglong4/',
    counts: {
      entity_clients: 171,
      non_entity_clients: 20,
      clients: 191,
    },
    mounts: [
      {
        mount_path: 'auth/method/uMGBU',
        "counts":{
          "distinct_entities":0,
          "entity_clients":0,
          "non_entity_tokens":0,
          "non_entity_clients":10,
          "clients":10
        }
      },
    ],
  },
];

transformedPayload.by_namespace = [
  {
    label: 'namespacelonglonglong4/',
    entity_clients: 171,
    non_entity_clients: 20,
    clients: 191,
    mounts: [
      {
        label: 'auth/method/uMGBU',
        entity_clients: 20,
        non_entity_clients: 15,
        clients: 35,
      },
    ],
  },
]
*/
