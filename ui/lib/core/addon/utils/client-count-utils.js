import { parseAPITimestamp } from 'core/utils/date-formatters';
import { compareAsc } from 'date-fns';

export const formatByMonths = (monthsArray) => {
  // the months array will always include a timestamp of the month and either new/total client data or counts = null
  if (!Array.isArray(monthsArray)) return monthsArray;

  const sortedPayload = sortMonthsByTimestamp(monthsArray);
  return sortedPayload?.map((m) => {
    const month = parseAPITimestamp(m.timestamp, 'M/yy');
    let totalClientsByNamespace = formatByNamespace(m.namespaces);
    let newClientsByNamespace = formatByNamespace(m.new_clients?.namespaces);
    if (Object.keys(m).includes('counts')) {
      let totalCounts = flattenDataset(m);
      let newCounts = m.new_clients ? flattenDataset(m.new_clients) : {};
      return {
        month,
        ...totalCounts,
        namespaces: formatByNamespace(m.namespaces) || [],
        namespaces_by_key: namespaceArrayToObject(totalClientsByNamespace, newClientsByNamespace, month),
        new_clients: {
          month,
          ...newCounts,
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

export const flattenDataset = (object) => {
  if (object?.counts) {
    let flattenedObject = {};
    Object.keys(object['counts']).forEach((key) => (flattenedObject[key] = object['counts'][key]));
    return homogenizeClientNaming(flattenedObject);
  }
  return object;
};

export const sortMonthsByTimestamp = (monthsArray) => {
  const sortedPayload = [...monthsArray];
  return sortedPayload.sort((a, b) =>
    compareAsc(parseAPITimestamp(a.timestamp), parseAPITimestamp(b.timestamp))
  );
};

export const namespaceArrayToObject = (totalClientsByNamespace, newClientsByNamespace, month) => {
  if (!totalClientsByNamespace) return {}; // return if no data for that month
  // all 'new_client' data resides within a separate key of each month (see data structure below)
  // FIRST: iterate and nest respective 'new_clients' data within each namespace and mount object
  // note: this is happening within the month object
  const nestNewClientsWithinNamespace = totalClientsByNamespace?.map((ns) => {
    let newNamespaceCounts = newClientsByNamespace?.find((n) => n.label === ns.label);
    if (newNamespaceCounts) {
      let { label, clients, entity_clients, non_entity_clients } = newNamespaceCounts;
      let newClientsByMount = [...newNamespaceCounts?.mounts];
      let nestNewClientsWithinMounts = ns.mounts?.map((mount) => {
        let new_clients = newClientsByMount?.find((m) => m.label === mount.label) || {};
        return {
          ...mount,
          new_clients,
        };
      });
      return {
        ...ns,
        new_clients: {
          label,
          clients,
          entity_clients,
          non_entity_clients,
        },
        mounts: [...nestNewClientsWithinMounts],
      };
    }
    return {
      ...ns,
      new_clients: {},
    };
  });
  // SECOND: create a new object (namespace_by_key) in which each namespace label is a key
  let namespaces_by_key = {};
  nestNewClientsWithinNamespace?.forEach((namespaceObject) => {
    // THIRD: make another object within the namespace where each mount label is a key
    let mounts_by_key = {};
    namespaceObject.mounts.forEach((mountObject) => {
      mounts_by_key[mountObject.label] = {
        month,
        ...mountObject,
        new_clients: { month, ...mountObject.new_clients },
      };
    });

    let { label, clients, entity_clients, non_entity_clients, new_clients } = namespaceObject;
    namespaces_by_key[label] = {
      month,
      clients,
      entity_clients,
      non_entity_clients,
      new_clients: { month, ...new_clients },
      mounts_by_key,
    };
  });
  return namespaces_by_key;
  // structure of object returned
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

// API RESPONSE STRUCTURE:
// data: {
//   ** by_namespace organized in descending order of client count number **
//   by_namespace: [
//     {
//       namespace_id: '96OwG',
//       namespace_path: 'test-ns/',
//       counts: {},
//       mounts: [{ mount_path: 'path-1', counts: {} }],
//     },
//   ],
//   ** months organized in ascending order of timestamps, oldest to most recent
//   months: [
//     {
//       timestamp: '2022-03-01T00:00:00Z',
//       counts: {},
//       namespaces: [
//         {
//           namespace_id: 'root',
//           namespace_path: '',
//           counts: {},
//           mounts: [{ mount_path: 'auth/up2/', counts: {} }],
//         },
//       ],
//       new_clients: {
//         counts: {},
//         namespaces: [
//           {
//             namespace_id: 'root',
//             namespace_path: '',
//             counts: {},
//             mounts: [{ mount_path: 'auth/up2/', counts: {} }],
//           },
//         ],
//       },
//     },
//     {
//       timestamp: '2022-04-01T00:00:00Z',
//       counts: {},
//       namespaces: [
//         {
//           namespace_id: 'root',
//           namespace_path: '',
//           counts: {},
//           mounts: [{ mount_path: 'auth/up2/', counts: {} }],
//         },
//       ],
//       new_clients: {
//         counts: {},
//         namespaces: [
//           {
//             namespace_id: 'root',
//             namespace_path: '',
//             counts: {},
//             mounts: [{ mount_path: 'auth/up2/', counts: {} }],
//           },
//         ],
//       },
//     },
//   ],
//   start_time: 'start timestamp string',
//   end_time: 'end timestamp string',
//   total: { clients: 300, non_entity_clients: 100, entity_clients: 400} ,
// }
