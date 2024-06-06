/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { parseAPITimestamp } from 'core/utils/date-formatters';
import { compareAsc, getUnixTime, isWithinInterval } from 'date-fns';

import type ClientsConfigModel from 'vault/models/clients/config';
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

export const setStartTimeQuery = (
  isEnterprise: boolean,
  config: ClientsConfigModel | Record<string, never>
) => {
  // CE versions have no license and so the start time defaults to "0001-01-01T00:00:00Z"
  if (isEnterprise && _hasConfig(config)) {
    return getUnixTime(config.billingStartTimestamp);
  }
  return null;
};

// METHODS FOR SERIALIZING ACTIVITY RESPONSE
export const formatDateObject = (dateObj: { monthIdx: number; year: number }, isEnd: boolean) => {
  const { year, monthIdx } = dateObj;
  // day=0 for Date.UTC() returns the last day of the month before
  // increase monthIdx by one to get last day of queried month
  const utc = isEnd ? Date.UTC(year, monthIdx + 1, 0) : Date.UTC(year, monthIdx, 1);
  return getUnixTime(utc);
};

export const formatByMonths = (
  monthsArray: (ActivityMonthBlock | EmptyActivityMonthBlock | NoNewClientsActivityMonthBlock)[]
) => {
  const sortedPayload = sortMonthsByTimestamp(monthsArray);
  return sortedPayload?.map((m) => {
    const month = parseAPITimestamp(m.timestamp, 'M/yy') as string;
    const { timestamp } = m;
    // counts are null if there is no monthly data
    if (m.counts) {
      const totalClientsByNamespace = formatByNamespace(m.namespaces);
      const newClientsByNamespace = formatByNamespace(m.new_clients?.namespaces);

      let newClients: ByMonthNewClients = { month, timestamp, namespaces: [] };
      if (m.new_clients?.counts) {
        newClients = {
          month,
          timestamp,
          ...destructureClientCounts(m?.new_clients?.counts),
          namespaces: formatByNamespace(m.new_clients?.namespaces),
        };
      }
      return {
        month,
        timestamp,
        ...destructureClientCounts(m.counts),
        namespaces: formatByNamespace(m.namespaces),
        namespaces_by_key: namespaceArrayToObject(
          totalClientsByNamespace,
          newClientsByNamespace,
          month,
          m.timestamp
        ),
        new_clients: newClients,
      };
    }
    // empty month
    return {
      month,
      timestamp,
      namespaces: [],
      namespaces_by_key: {},
      new_clients: { month, timestamp, namespaces: [] },
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

// In 1.10 'distinct_entities' changed to 'entity_clients' and 'non_entity_tokens' to 'non_entity_clients'
// these deprecated keys still exist on the response, so only return relevant keys here
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

export const sortMonthsByTimestamp = (
  monthsArray: (ActivityMonthBlock | EmptyActivityMonthBlock | NoNewClientsActivityMonthBlock)[]
) => {
  const sortedPayload = [...monthsArray];
  return sortedPayload.sort((a, b) =>
    compareAsc(parseAPITimestamp(a.timestamp) as Date, parseAPITimestamp(b.timestamp) as Date)
  );
};

export const namespaceArrayToObject = (
  monthTotals: ByNamespaceClients[],
  // technically this arg (monthNew) is the same type as above, just nested inside monthly new clients
  monthNew: ByMonthClients['new_clients']['namespaces'] | null,
  month: string,
  timestamp: string
) => {
  // namespaces_by_key is used to filter monthly activity data by namespace
  // it's an object in each month data block where the keys are namespace paths
  // and values include new and total client counts for that namespace in that month
  const namespaces_by_key = monthTotals.reduce((nsObject: { [key: string]: NamespaceByKey }, ns) => {
    const keyedNs: NamespaceByKey = {
      ...destructureClientCounts(ns),
      timestamp,
      month,
      mounts_by_key: {},
      new_clients: {
        month,
        timestamp,
        label: ns.label,
        mounts: [],
      },
    };
    const newNsClients = monthNew?.find((n) => n.label === ns.label);
    // mounts_by_key is is used to filter further in a namespace and get monthly activity by mount
    // it's an object inside the namespace block where the keys are mount paths
    // and the values include new and total client counts for that mount in that month
    keyedNs.mounts_by_key = ns.mounts.reduce(
      (mountObj: { [key: string]: MountByKey }, mount) => {
        const mountNewClients = newNsClients ? newNsClients.mounts.find((m) => m.label === mount.label) : {};
        mountObj[mount.label] = {
          ...mount,
          timestamp,
          month,
          new_clients: {
            timestamp,
            month,
            label: mount.label,
            ...mountNewClients,
          },
        };

        return mountObj;
      },
      {} as { [key: string]: MountByKey }
    );
    if (newNsClients) {
      keyedNs.new_clients = { month, timestamp, ...newNsClients };
    }
    nsObject[ns.label] = keyedNs;
    return nsObject;
  }, {});

  return namespaces_by_key;
};

// type guards for conditionals
function _hasConfig(model: ClientsConfigModel | object): model is ClientsConfigModel {
  if (!model) return false;
  return 'billingStartTimestamp' in model;
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
  namespaces_by_key: { [key: string]: NamespaceByKey };
  new_clients: ByMonthNewClients;
}

// clients numbers are only returned if month is of type ActivityMonthBlock
export interface ByMonthNewClients extends TotalClientsSometimes {
  month: string;
  timestamp: string;
  namespaces: ByNamespaceClients[];
}

export interface NamespaceByKey extends TotalClients {
  month: string;
  timestamp: string;
  mounts_by_key: { [key: string]: MountByKey };
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

export interface ActivityMonthBlock {
  timestamp: string; // YYYY-MM-01T00:00:00Z (always the first day of the month)
  counts: Counts;
  namespaces: NamespaceObject[];
  new_clients: {
    counts: Counts;
    namespaces: NamespaceObject[];
    timestamp: string;
  };
}

export interface NoNewClientsActivityMonthBlock {
  timestamp: string; // YYYY-MM-01T00:00:00Z (always the first day of the month)
  counts: Counts;
  namespaces: NamespaceObject[];
  new_clients: {
    counts: null;
    namespaces: null;
  };
}

export interface EmptyActivityMonthBlock {
  timestamp: string; // YYYY-MM-01T00:00:00Z (always the first day of the month)
  counts: null;
  namespaces: null;
  new_clients: null;
}

export interface Counts {
  acme_clients: number;
  clients: number;
  distinct_entities: number;
  entity_clients: number;
  non_entity_clients: number;
  non_entity_tokens: number;
  secret_syncs: number;
}
