/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { parseAPITimestamp } from 'core/utils/date-formatters';
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

// map to dropdowns for filtering client count tables
export enum ClientFilters {
  NAMESPACE = 'namespace_path',
  MOUNT_PATH = 'mount_path',
  MOUNT_TYPE = 'mount_type',
}

export type ClientFilterTypes = (typeof ClientFilters)[keyof typeof ClientFilters];

// returns array of VersionHistoryModels for noteworthy upgrades: 1.9, 1.10
// that occurred between timestamps (i.e. queried activity data)
export const filterVersionHistory = (
  versionHistory: ClientsVersionHistoryModel[],
  start: string,
  end: string
) => {
  if (versionHistory && start && end) {
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
    const { timestamp } = m;
    if (monthIsEmpty(m)) {
      // empty month
      return {
        timestamp,
        namespaces: [],
        new_clients: { timestamp, namespaces: [] },
      };
    }

    let newClients: ByMonthNewClients = { timestamp, namespaces: [] };
    if (monthWithAllCounts(m)) {
      newClients = {
        timestamp,
        ...destructureClientCounts(m?.new_clients.counts),
        namespaces: formatByNamespace(m.new_clients.namespaces),
      };
    }
    return {
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
    const nsLabel = ns.namespace_path === '' ? ns.namespace_id : ns.namespace_path;
    // data prior to adding mount granularity will still have a mounts array,
    // but the mount_path value will be "no mount accessor (pre-1.10 upgrade?)" (ref: vault/activity_log_util_common.go)
    // transform to an empty array for type consistency
    let mounts: MountClients[] | [] = [];
    if (Array.isArray(ns.mounts)) {
      mounts = ns.mounts.map((m) => ({
        label: m.mount_path,
        namespace_path: nsLabel,
        mount_path: m.mount_path,
        mount_type: m.mount_type,
        ...destructureClientCounts(m.counts),
      }));
    }
    return {
      label: nsLabel,
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

export const filterTableData = (
  data: MountClients[],
  filters: Record<ClientFilterTypes, string>
): MountClients[] => {
  // Return original data if no filters are specified
  if (!filters || Object.values(filters).every((v) => !v)) {
    return data;
  }

  return data.filter((datum) => {
    // Datum must satisfy every filter
    return Object.entries(filters).every(([filterKey, filterValue]) => {
      // If no filter is specified for that key, return true
      if (!filterValue) return true;
      // Otherwise only return true if the datum matches the filter
      return datum[filterKey as ClientFilterTypes] === filterValue;
    });
  });
};

export const flattenMounts = (namespaceArray: ByNamespaceClients[]) =>
  namespaceArray.map((n) => n.mounts).flat();

// TYPE GUARDS FOR CONDITIONALS
function monthIsEmpty(month: ActivityMonthBlock): month is ActivityMonthEmpty {
  return !month || month?.counts === null;
}

function monthWithAllCounts(month: ActivityMonthBlock): month is ActivityMonthStandard {
  return month?.counts !== null && month?.new_clients?.counts !== null;
}

export function filterIsSupported(f: string): f is ClientFilterTypes {
  return Object.values(ClientFilters).includes(f as ClientFilterTypes);
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
  mount_path: string;
  mount_type: string;
  namespace_path: string;
}

export interface ByMonthClients extends TotalClients {
  timestamp: string;
  namespaces: ByNamespaceClients[];
  new_clients: ByMonthNewClients;
}

export interface ByMonthNewClients extends TotalClientsSometimes {
  timestamp: string;
  namespaces: ByNamespaceClients[];
}

export interface NamespaceByKey extends TotalClients {
  timestamp: string;
  new_clients: NamespaceNewClients;
}

export interface NamespaceNewClients extends TotalClientsSometimes {
  timestamp: string;
  label: string;
  mounts: MountClients[];
}

export interface MountByKey extends TotalClients {
  timestamp: string;
  label: string;
  new_clients: MountNewClients;
}

export interface MountNewClients extends TotalClientsSometimes {
  timestamp: string;
  label: string;
}

// Serialized data from activity/export API
export interface ActivityExportData {
  client_id: string;
  client_type: string;
  namespace_id: string;
  namespace_path: string;
  mount_accessor: string;
  mount_type: string;
  mount_path: string;
  token_creation_time: string;
  client_first_used_time: string;
}
export interface EntityClients extends ActivityExportData {
  entity_name: string;
  entity_alias_name: string;
  local_entity_alias: boolean;
  policies: string[];
  entity_metadata: Record<string, any>;
  entity_alias_metadata: Record<string, any>;
  entity_alias_custom_metadata: Record<string, any>;
  entity_group_ids: string[];
}

// API RESPONSE SHAPE (prior to serialization)

export interface NamespaceObject {
  namespace_id: string;
  namespace_path: string;
  counts: Counts;
  mounts: { mount_path: string; counts: Counts; mount_type: string }[];
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
