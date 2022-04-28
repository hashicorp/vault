import { parseAPITimestamp } from 'core/utils/date-formatters';

export const flattenDataset = (object) => {
  let flattenedObject = {};
  Object.keys(object['counts']).forEach((key) => (flattenedObject[key] = object['counts'][key]));
  return homogenizeClientNaming(flattenedObject);
};

export const formatByMonths = (monthsArray) => {
  const sortedPayload = [...monthsArray];
  // months are always returned from the API: [mostRecent...oldestMonth]
  sortedPayload.reverse();
  return sortedPayload.map((m) => {
    if (Object.keys(m).includes('counts')) {
      let totalClients = flattenDataset(m);
      let newClients = flattenDataset(m.new_clients);
      m = {
        month: parseAPITimestamp(m.timestamp, 'M/yy'),
        ...totalClients,
        namespaces: formatByNamespace(m.namespaces),
        new_clients: {
          month: parseAPITimestamp(m.timestamp, 'M/yy'),
          ...newClients,
          namespaces: formatByNamespace(m.new_clients.namespaces),
        },
      };
      return nestCountsWithinNamespaceKey(m);
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
};

export const formatByNamespace = (namespaceArray) => {
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

export const nestCountsWithinNamespaceKey = (month) => {
  // create new key of `by_namespace_key` for month object
  month.by_namespace_key = {};
  if (month.namespaces) {
    month.namespaces.forEach((namespace) => {
      let { clients, entity_clients, non_entity_clients, mounts } = namespace;
      let new_clients = {};
      if (month.new_clients) {
        new_clients = month.new_clients.namespaces.find((n) => n.label === namespace.label) || {};
      }
      // create counts object with namespace label as key name
      month.by_namespace_key[namespace.label] = {
        clients,
        entity_clients,
        non_entity_clients,
        mounts,
        new_clients,
      };
      // TODO delete or keep new_clients.label within namespace key object?
      // delete month.by_namespace_key[namespace.label].new_clients.label
    });
  }
  return month;
};
