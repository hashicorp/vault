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
    // data prior to adding mount granularity will still have a mounts key,
    // but with the value: "no mount accessor (pre-1.10 upgrade?)" (ref: vault/activity_log_util_common.go)
    // transform to an empty array for type consistency
    let mounts = [];
    if (Array.isArray(ns.mounts)) {
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

  // namespaces_by_key is used to filter monthly activity data by namespace
  // it's an object in each month data block where the keys are namespace paths
  // and values include new and total client counts for that namespace in that month
  const namespaces_by_key = totalClientsByNamespace.reduce((nsObject, ns) => {
    const newNsClients = newClientsByNamespace?.find((n) => n.label === ns.label) || { month };

    // mounts_by_key is is used to filter further in a namespace and get monthly activity by mount
    // it's an object inside the namespace block where the keys are mount paths
    // and the values include new and total client counts for that mount in that month
    const mounts_by_key = ns.mounts.reduce((mountObj, mount) => {
      const newMountClients = newNsClients.mounts.find((m) => m.label === mount.label) || { month };
      mountObj[mount.label] = { ...mount, timestamp, month, new_clients: { month, ...newMountClients } };
      return mountObj;
    }, {});

    nsObject[ns.label] = { ...ns, timestamp, month, new_clients: { month, ...newNsClients }, mounts_by_key };
    // remove unnecessary data
    delete nsObject[ns.label].mounts;
    delete nsObject[ns.label].label;
    return nsObject;
  }, {});

  return namespaces_by_key;
};
