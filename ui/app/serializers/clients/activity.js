/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from '../application';
import { formatISO } from 'date-fns';
import { formatByMonths, formatByNamespace, destructureClientCounts } from 'core/utils/client-count-utils';
import timestamp from 'core/utils/timestamp';

// see tests/helpers/clients/client-count-helpers for sample API response (ACTIVITY_RESPONSE_STUB)
// and transformed by_namespace and by_month examples (SERIALIZED_ACTIVITY_RESPONSE)
export default class ActivitySerializer extends ApplicationSerializer {
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (payload.id === 'no-data') {
      return super.normalizeResponse(store, primaryModelClass, payload, id, requestType);
    }
    const response_timestamp = formatISO(timestamp.now());
    const transformedPayload = {
      ...payload,
      response_timestamp,
      by_namespace: formatByNamespace(payload.data.by_namespace),
      by_month: formatByMonths(payload.data.months),
      total: destructureClientCounts(payload.data.total),
    };
    delete payload.data.by_namespace;
    delete payload.data.months;
    delete payload.data.total;
    return super.normalizeResponse(store, primaryModelClass, transformedPayload, id, requestType);
  }
}
