// Common sub-types for client counts data
interface ChartTimestamp {
  month: string; // eg. 12/22
  timestamp: string; // ISO 8601
}

// Count and EmptyCount are mutually exclusive
// but that's hard to represent in an interface
// so for now we just have both
interface Count {
  clients?: number;
  entity_clients?: number;
  non_entity_clients?: number;
  secret_syncs?: number;
}
interface EmptyCount {
  count?: null;
}

interface MountData extends Count, EmptyCount {
  label: string; // eg 'auth/authid1'
}
interface NamespaceData extends Count, EmptyCount {
  label: string; // eg 'ns/foo'
  mounts: MountData[];
}
interface NewClients extends ChartTimestamp, Count, EmptyCount {
  namespaces: NamespaceData[];
}
interface NamespacesByKey {
  [key: string]: NamespaceData;
}
export interface SerializedChartData extends Count, EmptyCount, ChartTimestamp {
  namespaces: NamespaceData[];
  namespaces_by_key: NamespacesByKey;
  new_clients: NewClients;
  [key: string]: unknown;
}
