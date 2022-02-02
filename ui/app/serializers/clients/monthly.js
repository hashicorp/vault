import ApplicationSerializer from '../application';
import { formatISO } from 'date-fns';
export default class MonthlySerializer extends ApplicationSerializer {
  flattenDataset(byNamespaceArray) {
    let topTen = byNamespaceArray ? byNamespaceArray.slice(0, 10) : [];

    return topTen.map((ns) => {
      // 'namespace_path' is an empty string for root
      if (ns['namespace_id'] === 'root') ns['namespace_path'] = 'root';
      let label = ns['namespace_path'];
      let flattenedNs = {};
      // we don't want client counts nested within the 'counts' object for stacked charts
      Object.keys(ns['counts']).forEach((key) => (flattenedNs[key] = ns['counts'][key]));

      flattenedNs = this.homogenizeClientNaming(flattenedNs);

      // if mounts attribution unavailable, mounts will be undefined
      flattenedNs.mounts = ns.mounts?.map((mount) => {
        let flattenedMount = {};
        flattenedMount.label = mount['path'];
        Object.keys(mount['counts']).forEach((key) => (flattenedMount[key] = mount['counts'][key]));
        return flattenedMount;
      });

      return {
        label,
        ...flattenedNs,
      };
    });
  }

  // For 1.10 release naming changed from 'distinct_entities' to 'entity_clients' and
  // 'non_entity_tokens' to 'non_entity_clients'
  // accounting for deprecated API keys here and updating to latest nomenclature
  homogenizeClientNaming(object) {
    // TODO CMB check with API payload, latest draft includes both new and old key names
    // Add else to delete old key names IF correct ones exist?
    if (Object.keys(object).includes('distinct_entities', 'non_entity_tokens')) {
      let entity_clients = object.distinct_entities;
      let non_entity_clients = object.non_entity_tokens;
      let { clients } = object;
      return {
        clients,
        entity_clients,
        non_entity_clients,
      };
    }
  }

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    let response_timestamp = formatISO(new Date());
    let transformedPayload = {
      ...payload,
      response_timestamp,
      by_namespace: this.flattenDataset(payload.data.by_namespace),
      // nest within 'total' object to mimic /activity response shape
      total: this.homogenizeClientNaming(payload.data),
    };
    delete payload.data.by_namespace;
    return super.normalizeResponse(store, primaryModelClass, transformedPayload, id, requestType);
  }
}
