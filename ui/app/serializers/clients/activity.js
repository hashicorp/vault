import ApplicationSerializer from '../application';
import { formatISO } from 'date-fns';
import { parseAPITimestamp, parseRFC3339 } from 'core/utils/date-formatters';
export default class ActivitySerializer extends ApplicationSerializer {
  flattenDataset(object) {
    let flattenedObject = {};
    Object.keys(object['counts']).forEach((key) => (flattenedObject[key] = object['counts'][key]));
    return this.homogenizeClientNaming(flattenedObject);
  }

  formatByNamespace(namespaceArray) {
    return namespaceArray?.map((ns) => {
      // 'namespace_path' is an empty string for root
      if (ns['namespace_id'] === 'root') ns['namespace_path'] = 'root';
      let label = ns['namespace_path'];
      let flattenedNs = this.flattenDataset(ns);
      // if no mounts, mounts will be an empty array
      flattenedNs.mounts = [];
      if (ns?.mounts && ns.mounts.length > 0) {
        flattenedNs.mounts = ns.mounts.map((mount) => {
          return {
            label: mount['mount_path'],
            ...this.flattenDataset(mount),
          };
        });
      }
      return {
        label,
        ...flattenedNs,
      };
    });
  }

  formatByMonths(monthsArray) {
    const sortedPayload = [...monthsArray];
    // months are always returned from the API: [mostRecent...oldestMonth]
    sortedPayload.reverse();
    return sortedPayload.map((m) => {
      if (Object.keys(m).includes('counts')) {
        let totalClients = this.flattenDataset(m);
        let newClients = this.flattenDataset(m.new_clients);
        return {
          month: parseAPITimestamp(m.timestamp, 'M/yy'),
          ...totalClients,
          namespaces: this.formatByNamespace(m.namespaces),
          new_clients: {
            month: parseAPITimestamp(m.timestamp, 'M/yy'),
            ...newClients,
            namespaces: this.formatByNamespace(m.new_clients.namespaces),
          },
        };
      }
      // TODO CMB below is an assumption, need to test
      // if no monthly data (no counts key), month object will just contain a timestamp
      return {
        month: parseAPITimestamp(m.timestamp, 'M/yy'),
        new_clients: {
          month: parseAPITimestamp(m.timestamp, 'M/yy'),
        },
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
    return object;
  }

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (payload.id === 'no-data') {
      return super.normalizeResponse(store, primaryModelClass, payload, id, requestType);
    }
    let response_timestamp = formatISO(new Date());
    let transformedPayload = {
      ...payload,
      response_timestamp,
      by_namespace: this.formatByNamespace(payload.data.by_namespace),
      by_month: this.formatByMonths(payload.data.months),
      total: this.homogenizeClientNaming(payload.data.total),
      formatted_end_time: parseRFC3339(payload.data.end_time),
      formatted_start_time: parseRFC3339(payload.data.start_time),
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
