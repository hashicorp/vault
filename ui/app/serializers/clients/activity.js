import ApplicationSerializer from '../application';
import { formatISO } from 'date-fns';
export default class ActivitySerializer extends ApplicationSerializer {
  flattenDataset(byNamespaceArray) {
    return byNamespaceArray.map((ns) => {
      // 'namespace_path' is an empty string for root
      if (ns['namespace_id'] === 'root') ns['namespace_path'] = 'root';
      let label = ns['namespace_path'];
      let flattenedNs = {};
      // we don't want client counts nested within the 'counts' object for stacked charts
      Object.keys(ns['counts']).forEach((key) => (flattenedNs[key] = ns['counts'][key]));
      flattenedNs = this.homogenizeClientNaming(flattenedNs);

      // TODO CMB check how this works with actual API endpoint
      // if no mounts, mounts will be an empty array
      flattenedNs.mounts = ns.mounts
        ? ns.mounts.map((mount) => {
            let flattenedMount = {};
            flattenedMount.label = mount['mount_path'];
            Object.keys(mount['counts']).forEach((key) => (flattenedMount[key] = mount['counts'][key]));
            return flattenedMount;
          })
        : [];

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
    // TODO CMB Delete old key names IF correct ones exist?
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
    return object;
  }

  parseRFC3339(timestamp) {
    // convert '2021-03-21T00:00:00Z' --> ['2021', 2] (e.g. 2021 March, month is zero indexed)
    return timestamp
      ? [timestamp.split('-')[0], Number(timestamp.split('-')[1].replace(/^0+/, '')) - 1]
      : null;
  }

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (payload.id === 'no-data') {
      return super.normalizeResponse(store, primaryModelClass, payload, id, requestType);
    }
    let response_timestamp = formatISO(new Date());
    let transformedPayload = {
      ...payload,
      response_timestamp,
      by_namespace: this.flattenDataset(payload.data.by_namespace),
      total: this.homogenizeClientNaming(payload.data.total),
      formatted_end_time: this.parseRFC3339(payload.data.end_time),
      formatted_start_time: this.parseRFC3339(payload.data.start_time),
    };
    delete payload.data.by_namespace;
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
