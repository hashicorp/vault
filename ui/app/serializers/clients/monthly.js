import ApplicationSerializer from '../application';
import { formatISO } from 'date-fns';
import { formatByNamespace, homogenizeClientNaming } from 'core/utils/client-count-utils';

export default class MonthlySerializer extends ApplicationSerializer {
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (payload.id === 'no-data') {
      return super.normalizeResponse(store, primaryModelClass, payload, id, requestType);
    }
    const response_timestamp = formatISO(new Date());
    // TODO CMB: the following is assumed, need to confirm
    // the months array will always include a single object: a timestamp of the current month and new/total count data, if available
    const transformedPayload = {
      ...payload,
      response_timestamp,
      by_namespace: formatByNamespace(payload.data.by_namespace),
      // nest within 'total' object to mimic /activity response shape
      total: homogenizeClientNaming(payload.data),
    };
    delete payload.data.by_namespace;
    delete payload.data.months;
    delete payload.data.total;
    return super.normalizeResponse(store, primaryModelClass, transformedPayload, id, requestType);
  }
}
