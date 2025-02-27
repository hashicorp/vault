/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import isEmpty from '@ember/utils/lib/is_empty';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { sanitizePath } from 'core/utils/sanitize-path';
import { compareAsc, getUnixTime, isWithinInterval } from 'date-fns';
import type ClientsVersionHistoryModel from 'vault/vault/models/clients/version-history';

/*
The client count utils are responsible for serializing the sys/internal/counters/activity API response
The initial API response shape and serialized types are defined below.

To help visualize there are sample responses in ui/tests/helpers/clients.js
*/

// add new types here
export const CLIENT_TYPES = [
  'acme_clients',
  'clients', // summation of total clients
  'entity_clients',
  'non_entity_clients',
  'secret_syncs',
] as const;

export type ClientTypes = (typeof CLIENT_TYPES)[number];

// generates a block of total clients with 0's for use as defaults
function emptyCounts() {
  return CLIENT_TYPES.reduce((prev, type) => {
    const key = type;
    prev[key as ClientTypes] = 0;
    return prev;
  }, {} as TotalClientsSometimes) as TotalClients;
}

// returns array of VersionHistoryModels for noteworthy upgrades: 1.9, 1.10
// that occurred between timestamps (i.e. queried activity data)
export const filterVersionHistory = (
  versionHistory: ClientsVersionHistoryModel[],
  start: string,
  end: string
) => {
  if (versionHistory) {
    const upgrades = versionHistory.reduce((array: ClientsVersionHistoryModel[], upgradeData) => {
      const isRelevantHistory = (v: string) => {
        return (
          upgradeData.version.match(v) &&
          // only add if there is a previous version, otherwise this upgrade is the users' first version
          upgradeData.previousVersion &&
          // only add first match, disregard subsequent patch releases of the same version
          !array.some((d: ClientsVersionHistoryModel) => d.version.match(v))
        );
      };

      ['1.9', '1.10', '1.17'].forEach((v) => {
        if (isRelevantHistory(v)) array.push(upgradeData);
      });

      return array;
    }, []);

    // if there are noteworthy upgrades, only return those during queried date range
    if (upgrades.length) {
      const startDate = parseAPITimestamp(start) as Date;
      const endDate = parseAPITimestamp(end) as Date;
      return upgrades.filter(({ timestampInstalled }) => {
        const upgradeDate = parseAPITimestamp(timestampInstalled) as Date;
        return isWithinInterval(upgradeDate, { start: startDate, end: endDate });
      });
    }
  }
  return [];
};

// This method is used to return totals relevant only to the specified
// mount path within the specified namespace.
export const filteredTotalForMount = (
  byNamespace: ByNamespaceClients[],
  nsPath: string,
  mountPath: string
): TotalClients => {
  if (!nsPath || !mountPath || isEmpty(byNamespace)) return emptyCounts();
  return (
    byNamespace
      .find((namespace) => sanitizePath(namespace.label) === sanitizePath(nsPath))
      ?.mounts.find((mount: MountClients) => sanitizePath(mount.label) === sanitizePath(mountPath)) ||
    emptyCounts()
  );
};

// This method is used to filter byMonth data and return data for only
// the specified mount within the specified namespace. If data exists
// for the month but not the mount, it should return zero'd data. If
// no data exists for the month is returns the month as-is.
export const filterByMonthDataForMount = (
  byMonth: ByMonthClients[],
  namespacePath: string,
  mountPath: string
): ByMonthClients[] => {
  if (byMonth && namespacePath && mountPath) {
    const months: ByMonthClients[] = JSON.parse(JSON.stringify(byMonth));
    return [...months].map((m) => {
      if (m?.clients === undefined) {
        // if the month doesn't have data we can just return the block
        return m;
      }

      const nsData = m.namespaces?.find((ns) => sanitizePath(ns.label) === sanitizePath(namespacePath));
      const mountData = nsData?.mounts.find((mount) => sanitizePath(mount.label) === sanitizePath(mountPath));
      if (mountData) {
        // if we do have mount data, we need to add in new_client namespace information
        const nsNew = m.new_clients?.namespaces?.find(
          (ns) => sanitizePath(ns.label) === sanitizePath(namespacePath)
        );
        const mountNew =
          nsNew?.mounts.find((mount) => sanitizePath(mount.label) === sanitizePath(mountPath)) ||
          emptyCounts();
        return {
          month: m.month,
          timestamp: m.timestamp,
          ...mountData,
          namespaces: [], // this is just for making TS happy, matching the ByMonthClients shape
          new_clients: {
            month: m.month,
            timestamp: m.timestamp,
            label: mountPath,
            namespaces: [], // this is just for making TS happy, matching the ByMonthClients shape
            ...mountNew,
          },
        } as ByMonthClients;
      }
      // if the month has data but none for this mount, return mocked zeros
      return {
        month: m.month,
        timestamp: m.timestamp,
        label: mountPath,
        namespaces: [], // this is just for making TS happy, matching the ByMonthClients shape
        ...emptyCounts(),
        new_clients: {
          timestamp: m.timestamp,
          month: m.month,
          label: mountPath,
          namespaces: [], // this is just for making TS happy, matching the ByMonthClients shape
          ...emptyCounts(),
        },
      } as ByMonthClients;
    });
  }
  return byMonth;
};

// METHODS FOR SERIALIZING ACTIVITY RESPONSE
export const formatDateObject = (dateObj: { monthIdx: number; year: number }, isEnd: boolean) => {
  const { year, monthIdx } = dateObj;
  // day=0 for Date.UTC() returns the last day of the month before
  // increase monthIdx by one to get last day of queried month
  const utc = isEnd ? Date.UTC(year, monthIdx + 1, 0) : Date.UTC(year, monthIdx, 1);
  return getUnixTime(utc);
};

export const formatByMonths = (monthsArray: ActivityMonthBlock[]): ByMonthNewClients[] => {
  const sortedPayload = sortMonthsByTimestamp(monthsArray);
  return sortedPayload?.map((m) => {
    const month = parseAPITimestamp(m.timestamp, 'M/yy') as string;
    const { timestamp } = m;
    if (monthIsEmpty(m)) {
      // empty month
      return {
        month,
        timestamp,
        namespaces: [],
        new_clients: { month, timestamp, namespaces: [] },
      };
    }

    let newClients: ByMonthNewClients = { month, timestamp, namespaces: [] };
    if (monthWithAllCounts(m)) {
      newClients = {
        month,
        timestamp,
        ...destructureClientCounts(m?.new_clients.counts),
        namespaces: formatByNamespace(m.new_clients.namespaces),
      };
    }
    return {
      month,
      timestamp,
      ...destructureClientCounts(m.counts),
      namespaces: formatByNamespace(m.namespaces),
      new_clients: newClients,
    };
  });
};

export const formatByNamespace = (namespaceArray: NamespaceObject[] | null): ByNamespaceClients[] => {
  if (!Array.isArray(namespaceArray)) return [];
  return namespaceArray.map((ns) => {
    // i.e. 'namespace_path' is an empty string for 'root', so use namespace_id
    const label = ns.namespace_path === '' ? ns.namespace_id : ns.namespace_path;
    // data prior to adding mount granularity will still have a mounts array,
    // but the mount_path value will be "no mount accessor (pre-1.10 upgrade?)" (ref: vault/activity_log_util_common.go)
    // transform to an empty array for type consistency
    let mounts: MountClients[] | [] = [];
    if (Array.isArray(ns.mounts)) {
      mounts = ns.mounts.map((m) => ({ label: m.mount_path, ...destructureClientCounts(m.counts) }));
    }
    return {
      label,
      ...destructureClientCounts(ns.counts),
      mounts,
    };
  });
};

// This method returns only client types from the passed object, excluding other keys such as "label".
// when querying historical data the response will always contain the latest client type keys because the activity log is
// constructed based on the version of Vault the user is on (key values will be 0)
export const destructureClientCounts = (verboseObject: Counts | ByNamespaceClients) => {
  return CLIENT_TYPES.reduce(
    (newObj: Record<ClientTypes, Counts[ClientTypes]>, clientType: ClientTypes) => {
      newObj[clientType] = verboseObject[clientType];
      return newObj;
    },
    {} as Record<ClientTypes, Counts[ClientTypes]>
  );
};

export const sortMonthsByTimestamp = (monthsArray: ActivityMonthBlock[]) => {
  const sortedPayload = [...monthsArray];
  return sortedPayload.sort((a, b) =>
    compareAsc(parseAPITimestamp(a.timestamp) as Date, parseAPITimestamp(b.timestamp) as Date)
  );
};

// type guards for conditionals
function monthIsEmpty(month: ActivityMonthBlock): month is ActivityMonthEmpty {
  return !month || month?.counts === null;
}

function monthWithAllCounts(month: ActivityMonthBlock): month is ActivityMonthStandard {
  return month?.counts !== null && month?.new_clients?.counts !== null;
}

export function hasMountsKey(
  obj: ByMonthNewClients | NamespaceNewClients | MountNewClients
): obj is NamespaceNewClients {
  return 'mounts' in obj;
}

export function hasNamespacesKey(
  obj: ByMonthNewClients | NamespaceNewClients | MountNewClients
): obj is ByMonthNewClients {
  return 'namespaces' in obj;
}

// TYPES RETURNED BY UTILS (serialized)
export interface TotalClients {
  clients: number;
  entity_clients: number;
  non_entity_clients: number;
  secret_syncs: number;
  acme_clients: number;
}

// extend this type when the counts are optional (eg for new clients)
interface TotalClientsSometimes {
  clients?: number;
  entity_clients?: number;
  non_entity_clients?: number;
  secret_syncs?: number;
  acme_clients?: number;
}

export interface ByNamespaceClients extends TotalClients {
  label: string;
  mounts: MountClients[];
}

export interface MountClients extends TotalClients {
  label: string;
}

export interface ByMonthClients extends TotalClients {
  month: string;
  timestamp: string;
  namespaces: ByNamespaceClients[];
  new_clients: ByMonthNewClients;
}

export interface ByMonthNewClients extends TotalClientsSometimes {
  month: string;
  timestamp: string;
  namespaces: ByNamespaceClients[];
}

export interface NamespaceByKey extends TotalClients {
  month: string;
  timestamp: string;
  new_clients: NamespaceNewClients;
}

export interface NamespaceNewClients extends TotalClientsSometimes {
  month: string;
  timestamp: string;
  label: string;
  mounts: MountClients[];
}

export interface MountByKey extends TotalClients {
  month: string;
  timestamp: string;
  label: string;
  new_clients: MountNewClients;
}

export interface MountNewClients extends TotalClientsSometimes {
  month: string;
  timestamp: string;
  label: string;
}

// API RESPONSE SHAPE (prior to serialization)

export interface NamespaceObject {
  namespace_id: string;
  namespace_path: string;
  counts: Counts;
  mounts: { mount_path: string; counts: Counts }[];
}

type ActivityMonthStandard = {
  timestamp: string; // YYYY-MM-01T00:00:00Z (always the first day of the month)
  counts: Counts;
  namespaces: NamespaceObject[];
  new_clients: {
    counts: Counts;
    namespaces: NamespaceObject[];
    timestamp: string;
  };
};
type ActivityMonthNoNewClients = {
  timestamp: string; // YYYY-MM-01T00:00:00Z (always the first day of the month)
  counts: Counts;
  namespaces: NamespaceObject[];
  new_clients: {
    counts: null;
    namespaces: null;
  };
};
type ActivityMonthEmpty = {
  timestamp: string; // YYYY-MM-01T00:00:00Z (always the first day of the month)
  counts: null;
  namespaces: null;
  new_clients: null;
};
export type ActivityMonthBlock = ActivityMonthEmpty | ActivityMonthNoNewClients | ActivityMonthStandard;

export interface Counts {
  acme_clients: number;
  clients: number;
  entity_clients: number;
  non_entity_clients: number;
  secret_syncs: number;
}
