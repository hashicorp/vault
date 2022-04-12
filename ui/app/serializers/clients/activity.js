import ApplicationSerializer from '../application';
import { formatISO } from 'date-fns';
import { parseAPITimestamp, parseRFC3339 } from 'core/utils/date-formatters';
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

  // for vault usage - vertical bar chart
  flattenByMonths(payload, isNewClients = false) {
    const sortedPayload = [...payload];
    sortedPayload.reverse();
    if (isNewClients) {
      return sortedPayload?.map((m) => {
        return {
          month: parseAPITimestamp(m.timestamp, 'M/yy'),
          entity_clients: m.new_clients.counts.entity_clients,
          non_entity_clients: m.new_clients.counts.non_entity_clients,
          total: m.new_clients.counts.clients,
          namespaces: this.flattenDataset(m.new_clients.namespaces),
        };
      });
    } else {
      return sortedPayload?.map((m) => {
        return {
          month: parseAPITimestamp(m.timestamp, 'M/yy'),
          entity_clients: m.counts.entity_clients,
          non_entity_clients: m.counts.non_entity_clients,
          total: m.counts.clients,
          namespaces: this.flattenDataset(m.namespaces),
          new_clients: {
            entity_clients: m.new_clients.counts.entity_clients,
            non_entity_clients: m.new_clients.counts.non_entity_clients,
            total: m.new_clients.counts.clients,
            namespaces: this.flattenDataset(m.new_clients.namespaces),
          },
        };
      });
    }
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
    let response_timestamp = formatISO(new Date());
    let transformedPayload = {
      ...payload,
      response_timestamp,
      by_namespace: this.flattenDataset(payload.data.by_namespace),
      by_month_total_clients: this.flattenByMonths(payload.data.months),
      by_month_new_clients: this.flattenByMonths(payload.data.months, { isNewClients: true }),
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
