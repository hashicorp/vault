import { parseAPITimestamp } from 'core/utils/date-formatters';
import { compareAsc } from 'date-fns';

export const formatByMonths = (monthsArray) => {
  if (!Array.isArray(monthsArray)) return monthsArray;
  const sortedPayload = sortMonthsByTimestamp(monthsArray);
  return sortedPayload.map((m) => {
    const month = parseAPITimestamp(m.timestamp, 'M/yy');
    const totalClientsByNamespace = formatByNamespace(m.namespaces);
    const newClientsByNamespace = formatByNamespace(m.new_clients?.namespaces);
    if (Object.keys(m).includes('counts')) {
      let totalCounts = flattenDataset(m);
      let newCounts = m.new_clients ? flattenDataset(m.new_clients) : {};
      return {
        month,
        ...totalCounts,
        namespaces_by_key: namespaceArrayToObject(totalClientsByNamespace, newClientsByNamespace, month),
        new_clients: {
          month,
          ...newCounts,
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

const flattenDataset = (object) => {
  if (Object.keys(object).includes('counts') && object.counts) {
    let flattenedObject = {};
    Object.keys(object['counts']).forEach((key) => (flattenedObject[key] = object['counts'][key]));
    return homogenizeClientNaming(flattenedObject);
  }
  return object;
};

const sortMonthsByTimestamp = (monthsArray) => {
  // backend is working on a fix to sort months by date
  // right now months are ordered in descending client count number
  const sortedPayload = [...monthsArray];
  return sortedPayload.sort((a, b) =>
    compareAsc(parseAPITimestamp(a.timestamp), parseAPITimestamp(b.timestamp))
  );
};

const namespaceArrayToObject = (totalClientsByNamespace, newClientsByNamespace, month) => {
  // first nest 'new_clients' within each namespace and mount object respectively
  totalClientsByNamespace.forEach((ns) => {
    const newNamespaceCounts = newClientsByNamespace?.find((n) => n.label === ns.label);
    const newClientsByMount = newNamespaceCounts?.mounts;
    if (newClientsByMount) delete newNamespaceCounts.mounts;

    ns.new_clients = newNamespaceCounts || {};
    ns.mounts.forEach((mount) => {
      let newMountCounts = newClientsByMount?.find((m) => m.label === mount.label);
      mount.new_clients = newMountCounts || {};
    });
  });

  // create a new object in which each namespace label is a key
  // add month data to each object, including the nested 'new_clients', for time series graphs
  let namespaces_by_key = {};
  totalClientsByNamespace.forEach((namespaceObject) => {
    // and another where each mount label is a key
    let mounts_by_key = {};
    namespaceObject.mounts.forEach((mountObject) => {
      if (mountObject.new_clients) mountObject.new_clients.month = month;
      mounts_by_key[mountObject.label] = {
        month,
        ...mountObject,
      };
      // TODO CMB ^ delete the label key from final object don't think it's necessary
    });
    if (namespaceObject.new_clients) namespaceObject.new_clients.month = month;
    if (namespaceObject.mounts) delete namespaceObject.mounts; // the 'mounts_by_key' object replaces the 'mounts' array
    namespaces_by_key[namespaceObject.label] = {
      month,
      ...namespaceObject,
      mounts_by_key,
    };
  });
  return namespaces_by_key;
  // sample structure of object
  // namespace_by_key: {
  //   "namespace_label": {
  //     month: "3/22",
  //     clients: 32,
  //     entity_clients: 16,
  //     non_entity_clients: 16,
  //     new_clients: {
  //       month: "3/22",
  //       clients: 5,
  //       entity_clients: 2,
  //       non_entity_clients: 3,
  //     },
  //     mounts_by_key: {
  //       "mount_label": {
  //          month: "3/22",
  //          clients: 3,
  //          entity_clients: 2,
  //          non_entity_clients: 1,
  //          new_clients: {
  //           month: "3/22",
  //           clients: 5,
  //           entity_clients: 2,
  //           non_entity_clients: 3,
  //         },
  //       },
  //     },
  //   },
  // };
};
