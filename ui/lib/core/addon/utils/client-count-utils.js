/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { parseAPITimestamp } from 'core/utils/date-formatters';
import { compareAsc, getUnixTime, isWithinInterval } from 'date-fns';

// add new types here
export const CLIENT_TYPES = [
  'acme_clients',
  'clients', // summation of total clients
  'entity_clients',
  'non_entity_clients',
  'secret_syncs',
];

// returns array of VersionHistoryModels for noteworthy upgrades: 1.9, 1.10
// that occurred between timestamps (i.e. queried activity data)
export const filterVersionHistory = (versionHistory, start, end) => {
  if (versionHistory) {
    const upgrades = versionHistory.reduce((array, upgradeData) => {
      const includesVersion = (v) =>
        // only add first match, disregard subsequent patch releases of the same version
        upgradeData.version.match(v) && !array.some((d) => d.version.match(v));

      ['1.9', '1.10'].forEach((v) => {
        if (includesVersion(v)) array.push(upgradeData);
      });

      return array;
    }, []);

    // if there are noteworthy upgrades, only return those during queried date range
    if (upgrades.length) {
      const startDate = parseAPITimestamp(start);
      const endDate = parseAPITimestamp(end);
      return upgrades.filter(({ timestampInstalled }) => {
        const upgradeDate = parseAPITimestamp(timestampInstalled);
        return isWithinInterval(upgradeDate, { start: startDate, end: endDate });
      });
    }
  }
  return [];
};

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
    return {
      month,
      timestamp: m.timestamp,
      ...destructureClientCounts(m?.counts),
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
        ...destructureClientCounts(m?.new_clients?.counts),
        namespaces: formatByNamespace(m.new_clients?.namespaces) || [],
      },
    };
  });
};

export const formatByNamespace = (namespaceArray) => {
  if (!Array.isArray(namespaceArray)) return namespaceArray;
  return namespaceArray?.map((ns) => {
    // i.e. 'namespace_path' is an empty string for 'root', so use namespace_id
    const label = ns.namespace_path === '' ? ns.namespace_id : ns.namespace_path;
    // TODO ask backend what pre 1.10 data looks like, does "mounts" key exist?
    // if no mounts, mounts will be an empty array
    let mounts = [];
    if (ns?.mounts && ns.mounts.length > 0) {
      mounts = ns.mounts.map((m) => ({ label: m['mount_path'], ...destructureClientCounts(m?.counts) }));
    }
    return {
      label,
      ...destructureClientCounts(ns.counts),
      mounts,
    };
  });
};

// In 1.10 'distinct_entities' changed to 'entity_clients' and 'non_entity_tokens' to 'non_entity_clients'
// these deprecated keys still exist on the response, so only return relevant keys here
// when querying historical data the response will always contain the latest client type keys because the activity log is
// constructed based on the version of Vault the user is on (key values will be 0)
export const destructureClientCounts = (verboseObject) => {
  if (!verboseObject) return;
  return CLIENT_TYPES.reduce((newObj, clientType) => {
    newObj[clientType] = verboseObject[clientType];
    return newObj;
  }, {});
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
          label: ns.label,
          ...destructureClientCounts(newNamespaceCounts),
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

    const { label, new_clients } = namespaceObject;
    namespaces_by_key[label] = {
      month,
      timestamp,
      ...destructureClientCounts(namespaceObject),
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
