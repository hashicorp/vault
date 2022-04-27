import ApplicationSerializer from '../application';
import { formatISO } from 'date-fns';
import { formatByNamespace, homogenizeClientNaming } from 'core/utils/client-count-utils';

export default class MonthlySerializer extends ApplicationSerializer {
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (payload.id === 'no-data') {
      return super.normalizeResponse(store, primaryModelClass, payload, id, requestType);
    }
    let response_timestamp = formatISO(new Date());
    // TODO CMB: the following is assumed, need to confirm
    // the months array will always include a single object: a timestamp of the current month and new/total count data, if available
    let newClientsData = payload.data.months[0]?.new_clients || null;
    let by_namespace_new_clients, new_clients;
    if (newClientsData) {
      by_namespace_new_clients = formatByNamespace(newClientsData.namespaces);
      new_clients = homogenizeClientNaming(newClientsData.counts);
    } else {
      by_namespace_new_clients = [];
      new_clients = [];
    }
    let transformedPayload = {
      ...payload,
      response_timestamp,
      by_namespace_total_clients: formatByNamespace(payload.data.by_namespace),
      by_namespace_new_clients,
      // nest within 'total' object to mimic /activity response shape
      total: homogenizeClientNaming(payload.data),
      new: new_clients,
    };
    delete payload.data.by_namespace;
    delete payload.data.months;
    delete payload.data.total;
    return super.normalizeResponse(store, primaryModelClass, transformedPayload, id, requestType);
  }
}
