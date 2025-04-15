/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
type ISODateString = `${number}${number}-${number}${number}-${number}${number}`;
type ISOTimeString = `${number}${number}:${number}${number}:${number}${number}`;
type ISODateTimeString = `${ISODateString}T${ISOTimeString}`;
export interface TimeSeriesDatum {
    date: ISODateTimeString;
    value: number;
}
export interface SimpleDatum {
    value: number;
    label: string;
}
export declare enum REPLICATION_ENABLED_STATE {
    PRIMARY = "primary",
    SECONDARY = "secondary",
    BOOTSTRAPPING = "bootstrapping"
}
export declare const REPLICATION_DISABLED_STATE = "disabled";
export interface UsageDashboardData {
    auth_methods: Record<string, number>;
    kvv1_secrets: number;
    kvv2_secrets: number;
    lease_count_quotas: {
        global_lease_count_quota: {
            capacity: number;
            count: number;
            name: string;
        };
        total_lease_count_quotas: number;
    };
    namespaces: number;
    secrets_sync: number;
    pki: {
        total_issuers: number;
        total_roles: number;
    };
    replication_status: {
        dr_primary: boolean;
        dr_state: REPLICATION_ENABLED_STATE | typeof REPLICATION_DISABLED_STATE;
        pr_primary: boolean;
        pr_state: REPLICATION_ENABLED_STATE | typeof REPLICATION_DISABLED_STATE;
    };
    secret_engines: Record<string, number>;
}
export type getUsageDataFunction = () => Promise<UsageDashboardData>;
export {};
//# sourceMappingURL=index.d.ts.map