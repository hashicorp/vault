/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isSameMonthUTC, parseAPITimestamp } from 'core/utils/date-formatters';
import { compareAsc, isWithinInterval } from 'date-fns';
import { ROOT_NAMESPACE } from 'vault/services/namespace';
import { sanitizePath } from './sanitize-path';

import type ClientsVersionHistoryModel from 'vault/vault/models/clients/version-history';
import type {
  ActivityExportData,
  ActivityMonthBlock,
  ActivityMonthEmpty,
  ActivityMonthStandard,
  ByMonthNewClients,
  ByNamespaceClients,
  ClientFilterTypes,
  ClientTypes,
  Counts,
  MountClients,
  MountNewClients,
  NamespaceNewClients,
  NamespaceObject,
} from 'vault/vault/client-counts/activity-api';

/*
The client count utils are responsible for serializing the sys/internal/counters/activity API response
The initial API response shape and serialized types are defined below.

To help visualize there are sample responses in ui/tests/helpers/clients.js
*/

// Add new sys/activity/counters client count types here
export const CLIENT_TYPES = [
  'acme_clients',
  'clients', // summation of total clients
  'entity_clients',
  'non_entity_clients',
  'secret_syncs',
] as const;

// map to dropdowns for filtering client count tables
export enum ClientFilters {
  NAMESPACE = 'namespace_path',
  MOUNT_PATH = 'mount_path',
  MOUNT_TYPE = 'mount_type',
  // this filter/query param does not map to a key in either API response and is handled ~special~
  MONTH = 'month',
}

// client_type in the exported activity data differs slightly from the types of client keys
// returned by sys/internal/counters/activity endpoint (:
export const EXPORT_CLIENT_TYPES = ['non-entity-token', 'pki-acme', 'secret-sync', 'entity'] as const;

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
        // sanitized so it matches activity export data because mount_type there does NOT have a trailing slash
        mount_type: sanitizePath(m.mount_type),
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

// *Performance note*
// The client dashboard renders dropdown lists that specify filters. When the user selects a dropdown item (filter)
// it updates the query param and this method is called to filter the data passed to the displayed table.
// This method is not doing anything computationally expensive so it should be fine for filtering up to 50K rows of data.
// If activity data (either the by_namespace list or rows of data in the activity export API) grow past that, then we
// will want to look at converting this to a restartable task or do something else :)
export function filterTableData(
  data: MountClients[] | ActivityExportData[],
  filters: Record<ClientFilterTypes, string>
): MountClients[] | ActivityExportData[] {
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
      return matchesFilter(datum, filterKey as ClientFilterTypes, filterValue);
    });
  }) as typeof data;
}

const matchesFilter = (
  datum: ActivityExportData | MountClients,
  filterKey: ClientFilterTypes,
  filterValue: string
) => {
  // Only ActivityExportData data is ever filtered by 'client_first_used_time' (not MountClients)
  if (filterKey === ClientFilters.MONTH) {
    return 'client_first_used_time' in datum
      ? isSameMonthUTC(datum.client_first_used_time, filterValue)
      : false;
  }

  const datumValue = datum[filterKey];
  // The API returns and empty string as the namespace_path for the "root" namespace.
  // When a user selects "root" as a namespace filter we need to match the datum value
  // as either an empty string (for the activity export data) OR as "root"
  // (the by_namespace data is serialized to make "root" the namespace_path).
  if (filterKey === ClientFilters.NAMESPACE && filterValue === 'root') {
    return datumValue === ROOT_NAMESPACE || datumValue === filterValue;
  }
  return datumValue === filterValue;
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
