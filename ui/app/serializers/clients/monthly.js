import ApplicationSerializer from '../application';
import { formatISO } from 'date-fns';

export default class MonthlySerializer extends ApplicationSerializer {
  flattenDataset(namespaceArray) {
    return namespaceArray.map((ns) => {
      // 'namespace_path' is an empty string for root
      if (ns['namespace_id'] === 'root') ns['namespace_path'] = 'root';
      let label = ns['namespace_path'];
      let flattenedNs = {};
      // we don't want client counts nested within the 'counts' object for stacked charts
      Object.keys(ns['counts']).forEach((key) => (flattenedNs[key] = ns['counts'][key]));
      flattenedNs = this.homogenizeClientNaming(flattenedNs);

      // if no mounts, mounts will be an empty array
      flattenedNs.mounts = ns.mounts
        ? ns.mounts.map((mount) => {
            let flattenedMount = {};
            let label = mount['mount_path'];
            Object.keys(mount['counts']).forEach((key) => (flattenedMount[key] = mount['counts'][key]));
            flattenedMount = this.homogenizeClientNaming(flattenedMount);
            return {
              label,
              ...flattenedMount,
            };
          })
        : [];

      return {
        label,
        ...flattenedNs,
      };
    });
  }

  // In 1.10 'distinct_entities' changed to 'entity_clients' and
  // 'non_entity_tokens' to 'non_entity_clients'
  homogenizeClientNaming(object) {
    // if new key names exist, only return those key/value pairs
    if (Object.keys(object).includes('entity_clients')) {
      let { clients, entity_clients, non_entity_clients } = object;
      return {
        clients,
        entity_clients,
        non_entity_clients,
      };
    }
    // if object only has outdated key names, update naming
    if (Object.keys(object).includes('distinct_entities')) {
      let { clients, distinct_entities, non_entity_tokens } = object;
      return {
        clients,
        entity_clients: distinct_entities,
        non_entity_clients: non_entity_tokens,
      };
    }
    // TODO CMB: test what to return if neither key exists
    return object;
  }

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (payload.id === 'no-data') {
      return super.normalizeResponse(store, primaryModelClass, payload, id, requestType);
    }
    // TODO CMB will there always be a months key on this response?
    let newClientsData = payload.data.months[0]?.new_clients;
    let response_timestamp = formatISO(new Date());
    let transformedPayload = {
      ...payload,
      response_timestamp,
      by_namespace_total_clients: this.flattenDataset(payload.data.by_namespace),
      by_namespace_new_clients: this.flattenDataset(newClientsData.namespaces),
      // nest within 'total' object to mimic /activity response shape
      total: this.homogenizeClientNaming(payload.data),
      new: this.homogenizeClientNaming(newClientsData.counts),
    };
    delete payload.data.by_namespace;
    delete payload.data.months;
    delete payload.data.total;
    return super.normalizeResponse(store, primaryModelClass, transformedPayload, id, requestType);
  }
}
