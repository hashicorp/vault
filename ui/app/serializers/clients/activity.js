import ApplicationSerializer from '../application';
import { formatISO } from 'date-fns';

export default ApplicationSerializer.extend({
  flattenDataset(payload) {
    // let topTen = payload ? payload.slice(0, 10) : [];
    return payload?.map((ns) => {
      // 'namespace_path' is an empty string for root
      if (ns['namespace_id'] === 'root') ns['namespace_path'] = 'root';
      let label = ns['namespace_path'] || ns['namespace_id'];
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
    let response_timestamp = formatISO(new Date());
    let transformedPayload = {
      ...payload,
      response_timestamp,
      by_namespace: this.flattenDataset(payload.data.by_namespace),
    };
    delete payload.data.by_namespace;
    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },
});

/* 
SAMPLE PAYLOAD BEFORE/AFTER:

payload.data.by_namespace = [
  {
    namespace_id: '5SWT8',
    namespace_path: 'namespacelonglonglong4/',
    _comment1: 'client counts are nested within own object', 
    counts: {
      entity_clients: 171,
      non_entity_clients: 20,
      clients: 191,
    },
    mounts: [
      {
        path: 'auth/method/uMGBU',
        counts: {
          clients: 35,
          entity_clients: 20,
          non_entity_clients: 15,
        },
      },
    ],
  },
];

transformedPayload.by_namespace = [
  {
    label: 'namespacelonglonglong4/',
    _comment2: 'remove nested object', 
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
