/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { parseAPITimestamp } from 'core/utils/date-formatters';
import { compareAsc, getUnixTime } from 'date-fns';

export const formatDateObject = (dateObj, isEnd) => {
  if (dateObj) {
    const { year, monthIdx } = dateObj;
    // day=0 for Date.UTC() returns the last day of the month before
    // increase monthIdx by one to get last day of queried month
    const utc = isEnd ? Date.UTC(year, monthIdx + 1, 0) : Date.UTC(year, monthIdx, 1);
    return getUnixTime(utc);
  }
};

export const formatByMonths = (monthsArray) => {
  // the monthsArray will always include a timestamp of the month and either new/total client data or counts = null
  if (!Array.isArray(monthsArray)) return monthsArray;

  const sortedPayload = sortMonthsByTimestamp(monthsArray);
  return sortedPayload?.map((m) => {
    const month = parseAPITimestamp(m.timestamp, 'M/yy');
    const totalClientsByNamespace = formatByNamespace(m.namespaces);
    const newClientsByNamespace = formatByNamespace(m.new_clients?.namespaces);
    if (Object.keys(m).includes('counts')) {
      const totalCounts = flattenDataset(m);
      const newCounts = m.new_clients ? flattenDataset(m.new_clients) : {};
      return {
        month,
        timestamp: m.timestamp,
        ...totalCounts,
        namespaces: formatByNamespace(m.namespaces) || [],
        namespaces_by_key: namespaceArrayToObject(
          totalClientsByNamespace,
          newClientsByNamespace,
          month,
          m.timestamp
        ),
        new_clients: {
          month,
          timestamp: m.timestamp,
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
    const label = ns['namespace_path'];
    const flattenedNs = flattenDataset(ns);
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
    const { clients, entity_clients, non_entity_clients, secret_syncs } = object;
    return {
      clients,
      entity_clients,
      non_entity_clients,
      secret_syncs,
    };
  }
  // if object only has outdated key names, update naming
  if (Object.keys(object).includes('distinct_entities')) {
    const { clients, distinct_entities, non_entity_tokens } = object;
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
    const flattenedObject = {};
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

export const namespaceArrayToObject = (totalClientsByNamespace, newClientsByNamespace, month, timestamp) => {
  if (!totalClientsByNamespace) return {}; // return if no data for that month
  // all 'new_client' data resides within a separate key of each month (see data structure below)
  // FIRST: iterate and nest respective 'new_clients' data within each namespace and mount object
  // note: this is happening within the month object
  const nestNewClientsWithinNamespace = totalClientsByNamespace?.map((ns) => {
    const newNamespaceCounts = newClientsByNamespace?.find((n) => n.label === ns.label);
    if (newNamespaceCounts) {
      const { label, clients, entity_clients, non_entity_clients, secret_syncs } = newNamespaceCounts;
      const newClientsByMount = [...newNamespaceCounts.mounts];
      const nestNewClientsWithinMounts = ns.mounts?.map((mount) => {
        const new_clients = newClientsByMount?.find((m) => m.label === mount.label) || {};
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
          secret_syncs,
          mounts: newClientsByMount,
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
  const namespaces_by_key = {};
  nestNewClientsWithinNamespace?.forEach((namespaceObject) => {
    // THIRD: make another object within the namespace where each mount label is a key
    const mounts_by_key = {};
    namespaceObject.mounts.forEach((mountObject) => {
      mounts_by_key[mountObject.label] = {
        month,
        timestamp,
        ...mountObject,
        new_clients: { month, ...mountObject.new_clients },
      };
    });

    const { label, clients, entity_clients, non_entity_clients, secret_syncs, new_clients } = namespaceObject;
    namespaces_by_key[label] = {
      month,
      timestamp,
      clients,
      entity_clients,
      non_entity_clients,
      secret_syncs,
      new_clients: { month, ...new_clients },
      mounts_by_key,
    };
  });
  return namespaces_by_key;
  /*
  structure of object returned
  namespace_by_key: {
    "namespace_label": {
      month: "3/22",
      clients: 32,
      entity_clients: 16,
      non_entity_clients: 16,
      new_clients: {
        month: "3/22",
        clients: 5,
        entity_clients: 2,
        non_entity_clients: 3,
        mounts: [...array of this namespace's mounts and their new client counts],
      },
      mounts_by_key: {
        "mount_label": {
           month: "3/22",
           clients: 3,
           entity_clients: 2,
           non_entity_clients: 1,
           new_clients: {
            month: "3/22",
            clients: 5,
            entity_clients: 2,
            non_entity_clients: 3,
          },
        },
      },
    },
  };
  */
};

/*
API RESPONSE STRUCTURE:
data: {
  ** by_namespace organized in descending order of client count number **
  by_namespace: [
    {
      namespace_id: '96OwG',
      namespace_path: 'test-ns/',
      counts: {},
      mounts: [{ mount_path: 'path-1', counts: {} }],
    },
  ],
  ** months organized in ascending order of timestamps, oldest to most recent
  months: [
    {
      timestamp: '2022-03-01T00:00:00Z',
      counts: {},
      namespaces: [
        {
          namespace_id: 'root',
          namespace_path: '',
          counts: {},
          mounts: [{ mount_path: 'auth/up2/', counts: {} }],
        },
      ],
      new_clients: {
        counts: {},
        namespaces: [
          {
            namespace_id: 'root',
            namespace_path: '',
            counts: {},
            mounts: [{ mount_path: 'auth/up2/', counts: {} }],
          },
        ],
      },
    },
    {
      timestamp: '2022-04-01T00:00:00Z',
      counts: {},
      namespaces: [
        {
          namespace_id: 'root',
          namespace_path: '',
          counts: {},
          mounts: [{ mount_path: 'auth/up2/', counts: {} }],
        },
      ],
      new_clients: {
        counts: {},
        namespaces: [
          {
            namespace_id: 'root',
            namespace_path: '',
            counts: {},
            mounts: [{ mount_path: 'auth/up2/', counts: {} }],
          },
        ],
      },
    },
  ],
  start_time: 'start timestamp string',
  end_time: 'end timestamp string',
  total: { clients: 300, non_entity_clients: 100, entity_clients: 400} ,
}
*/
