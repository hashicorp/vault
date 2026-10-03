/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { parseAPITimestamp } from 'core/utils/date-formatters';
import { compareAsc, isValid, parseJSON } from 'date-fns';
import { sanitizePath } from '../sanitize-path';

import type {
  ActivityMonthBlock,
  ActivityMonthEmpty,
  ActivityMonthStandard,
  ByMonthNewClients,
  ByNamespaceClients,
  ClientTypes,
  Counts,
  MountClients,
  NamespaceObject,
} from 'vault/vault/client-counts/activity-api';
import { CLIENT_TYPES } from './helpers';

/*
These client count utils are responsible for serializing the sys/internal/counters/activity API response.
To help visualize there are sample responses in ui/tests/helpers/clients.js
*/

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

export const formatExportData = async (resp: Response, { isDownload = false }) => {
  // The response from the export API is a ReadableStream
  const blob = await resp.blob();
  // If the user wants to download the export data just return the blob.
  if (isDownload) return blob;

  // Otherwise format to JSON to render dataset in a table.
  const jsonLines = await blob.text();
  const lines = jsonLines.trim().split('\n');
  return lines.map((line: string) => JSON.parse(line));
};

export const formatQueryParams = (query: { start_time?: string; end_time?: string } = {}) => {
  const { start_time, end_time } = query;
  const formattedQuery: Partial<Record<'start_time' | 'end_time', string>> = {};

  if (start_time && isValid(parseJSON(start_time))) {
    formattedQuery.start_time = start_time;
  }
  if (end_time && isValid(parseJSON(end_time))) {
    formattedQuery.end_time = end_time;
  }

  return formattedQuery;
};

export const sortMonthsByTimestamp = (monthsArray: ActivityMonthBlock[]) => {
  const sortedPayload = [...monthsArray];
  return sortedPayload.sort((a, b) =>
    compareAsc(parseAPITimestamp(a.timestamp) as Date, parseAPITimestamp(b.timestamp) as Date)
  );
};

// TYPE GUARDS FOR CONDITIONALS
function monthIsEmpty(month: ActivityMonthBlock): month is ActivityMonthEmpty {
  return !month || month?.counts === null;
}

function monthWithAllCounts(month: ActivityMonthBlock): month is ActivityMonthStandard {
  return month?.counts !== null && month?.new_clients?.counts !== null;
}
