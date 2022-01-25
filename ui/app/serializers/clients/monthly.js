import ApplicationSerializer from '../application';
import { formatISO } from 'date-fns';

export default ApplicationSerializer.extend({
  flattenDataset(payload) {
    let topTen = payload.slice(0, 10);

    return topTen.map((ns) => {
      // 'namespace_path' is an empty string for root
      if (ns['namespace_id'] === 'root') ns['namespace_path'] = 'root';
      let label = ns['namespace_path'] || ns['namespace_id']; // TODO CMB will namespace_path ever be empty?
      let flattenedNs = {};
      // we don't want client counts nested within the 'counts' object for stacked charts
      Object.keys(ns['counts']).forEach((key) => (flattenedNs[key] = ns['counts'][key]));

      // homogenize client naming for all namespaces
      if (Object.keys(flattenedNs).includes('distinct_entities', 'non_entity_tokens')) {
        flattenedNs.entity_clients = flattenedNs.distinct_entities;
        flattenedNs.non_entity_clients = flattenedNs.non_entity_tokens;
        delete flattenedNs.distinct_entities;
        delete flattenedNs.non_entity_tokens;
      }

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
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    let { data } = payload;
    let { clients, distinct_entities, non_entity_tokens } = data;
    let response_timestamp = formatISO(new Date());
    let transformedPayload = {
      ...payload,
      // TODO CMB should these be nested under "data" to go to model correctly?)
      response_timestamp,
      by_namespace: this.flattenDataset(data.by_namespace),
      // nest within 'total' object to mimic /activity response shape
      total: {
        clients,
        entityClients: distinct_entities,
        nonEntityClients: non_entity_tokens,
      },
    };
    delete payload.data.by_namespace;
    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },
});
