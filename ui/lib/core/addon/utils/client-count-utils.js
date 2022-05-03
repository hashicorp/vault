import { parseAPITimestamp } from 'core/utils/date-formatters';
import { compareAsc } from 'date-fns';

export const flattenDataset = (object) => {
  if (Object.keys(object).includes('counts') && object.counts) {
    let flattenedObject = {};
    Object.keys(object['counts']).forEach((key) => (flattenedObject[key] = object['counts'][key]));
    return homogenizeClientNaming(flattenedObject);
  }
  return object;
};

export const formatByMonths = (monthsArray) => {
  if (!Array.isArray(monthsArray)) return monthsArray;
  const sortedPayload = sortMonthsByTimestamp(monthsArray);
  return sortedPayload.map((m) => {
    if (Object.keys(m).includes('counts')) {
      let totalClients = flattenDataset(m);
      let newClients = m.new_clients ? flattenDataset(m.new_clients) : {};
      return {
        month: parseAPITimestamp(m.timestamp, 'M/yy'),
        ...totalClients,
        namespaces: formatByNamespace(m.namespaces),
        new_clients: {
          month: parseAPITimestamp(m.timestamp, 'M/yy'),
          ...newClients,
          namespaces: formatByNamespace(m.new_clients?.namespaces) || [],
        },
      };
    }
  });
};

export const formatByNamespace = (namespaceArray) => {
  if (!Array.isArray(namespaceArray)) return namespaceArray;
  return namespaceArray?.map((ns) => {
    // 'namespace_path' is an empty string for root
    if (ns['namespace_id'] === 'root') ns['namespace_path'] = 'root';
    let label = ns['namespace_path'];
    let flattenedNs = flattenDataset(ns);
    // if no mounts, mounts will be an empty array
    flattenedNs.mounts = [];
    if (ns?.mounts && ns.mounts.length > 0) {
      flattenedNs.mounts = ns.mounts.map((mount) => {
        return {
          label: mount['mount_path'],
          ...flattenDataset(mount),
        };
      });
    }
    return {
      label,
      ...flattenedNs,
    };
  });
};

// In 1.10 'distinct_entities' changed to 'entity_clients' and
// 'non_entity_tokens' to 'non_entity_clients'
export const homogenizeClientNaming = (object) => {
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
};

export const sortMonthsByTimestamp = (monthsArray) => {
  // backend is working on a fix to sort months by date
  // right now months are ordered in descending client count number
  const sortedPayload = [...monthsArray];
  return sortedPayload.sort((a, b) =>
    compareAsc(parseAPITimestamp(a.timestamp), parseAPITimestamp(b.timestamp))
  );
};
